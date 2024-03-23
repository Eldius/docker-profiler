package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"sync"
)

var (
	ClientBuildErr = errors.New("failed to create Docker client")
)

type Client struct {
	d *client.Client
}

func NewClient() (*Client, error) {
	apiClient, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		err := fmt.Errorf("%w: %w", ClientBuildErr, err)
		return nil, err
	}
	return &Client{d: apiClient}, nil
}

func (c Client) GetRuntimeStatistcs(ctx context.Context) error {
	containerList, err := c.d.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	var finishFuncs []func()
	var wg sync.WaitGroup

	for _, instance := range containerList {
		wg.Add(1)
		fmt.Printf("- %v\n\n", instance.Names[0])
		s, err := c.d.ContainerStats(ctx, instance.ID, true)
		if err != nil {
			err = fmt.Errorf("fetching container status for '%s': %w", err)
			return err
		}
		finishFuncs = append(finishFuncs, func() {
			_ = s.Body.Close()
		})

		go func(wg *sync.WaitGroup) {
			sc := bufio.NewScanner(s.Body)
			defer wg.Done()
			var stats ContainerStats

			for sc.Scan() {
				_ = json.Unmarshal(sc.Bytes(), &stats)
				fmt.Printf("---\n- cpu:\n  - total usage: %v\n  - online: %v\n", stats.CPUStats.CPUUsage.TotalUsage, stats.CPUStats.OnlineCpus)
				fmt.Printf("\n- memory:\n  - limit: %vmb\n  - usage: %vkb\n", stats.MemoryStats.Limit/(1024*1024), stats.MemoryStats.Usage/(1024))
			}
		}(&wg)
	}

	wg.Wait()
	defer func() {
		for _, ff := range finishFuncs {
			ff()
		}
	}()

	return nil
}
