package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/eldius/docker-profiler/internal/docker"
	"github.com/eldius/docker-profiler/internal/plot"
	"log"
	"time"
)

func main() {

	containerName := flag.String("container", "", "Container name to be profiled")
	profile := flag.Bool("profile", false, "Profile containers")
	plotChart := flag.Bool("plot", false, "Profile containers")

	flag.Parse()

	fmt.Println("containerName:", *containerName)
	if *containerName == "" {
		panic(errors.New("invalid container name"))
	}
	fmt.Println("profile:", *profile)

	c, err := docker.NewClient(*containerName)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	if *profile {
		ctx := context.Background()

		if err := c.GetRuntimeStatistcs(ctx, *containerName); err != nil {
			log.Fatalf("failed to get runtime statistics: %+v", err)
		}
	}

	fmt.Println("plot:", *plotChart)
	if *plotChart {
		list, err := c.List()
		if err != nil {
			log.Fatalf("failed to list datapoints: %v", err)
		}

		for id, d := range list {
			fmt.Println("---")
			fmt.Printf("id:           %06d\n", id)
			fmt.Printf("timestamp:    %s\n", d.Timestamp.Format(time.RFC3339))
			fmt.Printf("memory usage: %s\n", d.MemoryUsageStr())
			fmt.Printf("memory limit: %s\n", d.MemoryLimitStr())
			fmt.Printf("cpu percent:  %01.2f\n", d.CPUPercentage)
			fmt.Printf("cpu online:   %01.2f\n", d.CPUOnlineCount)
			fmt.Printf("cpu usage:    %01.2f\n", d.CPUUsage)
			fmt.Printf("timestamp:    %v\n", d.Timestamp)
			fmt.Println("")
		}

		//plot.PlotToFile(list)
		plot.Plot(list)

	}
}
