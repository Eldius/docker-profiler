package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/eldius/docker-profiler/internal/model"
	"github.com/eldius/docker-profiler/internal/persistence"
	"strings"
	"sync"
)

var (
	ClientBuildErr = errors.New("failed to create Docker client")
)

type Client struct {
	d *client.Client
	r *persistence.Repository
}

func NewClient(name string) (*Client, error) {
	apiClient, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		err := fmt.Errorf("%w: %w", ClientBuildErr, err)
		return nil, err
	}
	return &Client{
		d: apiClient,
		r: persistence.NewRepository(name),
	}, nil
}

func (c Client) GetRuntimeStatistcs(ctx context.Context, name string) error {
	name = normalizeName(name)
	containerList, err := c.d.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	var finishFuncs []func()
	var wg sync.WaitGroup

	for _, instance := range containerList {
		iName := normalizeName(instance.Names[0])
		if strings.EqualFold(iName, name) {
			wg.Add(1)
			fmt.Printf("- %v\n\n", iName)
			s, err := c.d.ContainerStats(ctx, instance.ID, true)
			if err != nil {
				err = fmt.Errorf("fetching container status for '%s': %w", iName, err)
				return err
			}
			finishFuncs = append(finishFuncs, func() {
				_ = s.Body.Close()
			})

			go func(wg *sync.WaitGroup) {
				sc := bufio.NewScanner(s.Body)
				defer wg.Done()
				var stats model.ContainerStats

				for sc.Scan() {
					_ = json.Unmarshal(sc.Bytes(), &stats)
					if err := c.r.Persist(stats); err != nil {
						err = fmt.Errorf("persisting container stats: %w", err)
						panic(err)
					}
					fmt.Printf("---\n- cpu:\n  - total usage: %v\n  - percent usage: %01.2f%%\n  - online: %v\n", stats.CPUStats.CPUUsage.TotalUsage, stats.CPUUsagePercentage(), stats.CPUStats.OnlineCPUs)
					fmt.Printf("\n- memory:\n  - limit: %s\n  - usage: %s\n", stats.MemoryLimitStr(), stats.MemoryUsageStr())
				}
			}(&wg)
		}
	}

	wg.Wait()
	defer func() {
		for _, ff := range finishFuncs {
			ff()
		}
	}()

	return c.r.Close()
}

func (c Client) List() ([]model.MetricsDatapoint, error) {
	list, err := c.r.List()
	if err != nil {
		err = fmt.Errorf("trying to list datapoints: %w", err)
		return nil, err
	}

	return list, nil
}

func normalizeName(name string) string {
	return strings.TrimLeft(name, "/")
}
