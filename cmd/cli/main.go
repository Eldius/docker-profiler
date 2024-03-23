package main

import (
	"context"
	"flag"
	"github.com/eldius/docker-profiler/internal/docker"
	"log"
)

func main() {
	flag.Parse()

	c, err := docker.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	if err := c.GetRuntimeStatistcs(ctx); err != nil {
		log.Fatalf("failed to get runtime statistics: %+v", err)
	}
}
