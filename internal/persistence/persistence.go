package persistence

import (
	"fmt"
	"github.com/eldius/docker-profiler/internal/model"
	"github.com/nakabonne/tstorage"
	"os"
	"path/filepath"
	"time"
)

const (
	memoryUsageMetricName   = "memory_usage"
	memoryLimitMetricName   = "memory_limit"
	cpuOnlineMetricName     = "cpu_online"
	cpuUsageMetricName      = "cpu_usage"
	cpuPercentageMetricName = "cpu_percentage"
)

func NewRepository(name string) *Repository {
	dir, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("getting working directory: %w", err)
		panic(err)
	}
	storage, _ := tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Seconds),
		tstorage.WithDataPath(filepath.Join(dir, ".data", name)),
	)
	return &Repository{db: storage}
}

type Repository struct {
	db tstorage.Storage
}

func (r *Repository) Persist(s model.ContainerStats) error {
	unixTimestamp := time.Now().Unix()
	return r.db.InsertRows([]tstorage.Row{
		{
			Metric:    memoryUsageMetricName,
			DataPoint: tstorage.DataPoint{Timestamp: unixTimestamp, Value: float64(s.MemoryStats.Usage)},
		},
		{
			Metric:    memoryLimitMetricName,
			DataPoint: tstorage.DataPoint{Timestamp: unixTimestamp, Value: float64(s.MemoryStats.Usage)},
		},
		{
			Metric:    cpuOnlineMetricName,
			DataPoint: tstorage.DataPoint{Timestamp: unixTimestamp, Value: float64(s.CPUStats.OnlineCPUs)},
		},
		{
			Metric:    cpuUsageMetricName,
			DataPoint: tstorage.DataPoint{Timestamp: unixTimestamp, Value: float64(s.CPUStats.CPUUsage.TotalUsage)},
		},
		{
			Metric:    cpuPercentageMetricName,
			DataPoint: tstorage.DataPoint{Timestamp: unixTimestamp, Value: s.CPUUsagePercentage()},
		},
	})
}

func (r *Repository) List() ([]model.MetricsDatapoint, error) {
	mups, err := r.db.Select(memoryUsageMetricName, make([]tstorage.Label, 0), 0, time.Now().Unix())
	if err != nil {
		err = fmt.Errorf("listing memory usage datapoints: %w", err)
	}
	mlps, err := r.db.Select(memoryLimitMetricName, make([]tstorage.Label, 0), 0, time.Now().Unix())
	if err != nil {
		err = fmt.Errorf("listing memory limit datapoints: %w", err)
	}
	cops, err := r.db.Select(cpuOnlineMetricName, make([]tstorage.Label, 0), 0, time.Now().Unix())
	if err != nil {
		err = fmt.Errorf("listing online cpu count datapoints: %w", err)
	}
	cups, err := r.db.Select(cpuUsageMetricName, make([]tstorage.Label, 0), 0, time.Now().Unix())
	if err != nil {
		err = fmt.Errorf("listing online cpu count datapoints: %w", err)
	}
	cpps, err := r.db.Select(cpuPercentageMetricName, make([]tstorage.Label, 0), 0, time.Now().Unix())
	if err != nil {
		err = fmt.Errorf("listing online cpu count datapoints: %w", err)
	}

	resp := make([]model.MetricsDatapoint, len(mups))
	for i := range mups {
		resp[i] = model.MetricsDatapoint{
			Timestamp:      time.Unix(mups[i].Timestamp, 0),
			MemoryUsage:    getValue(mups, i),
			MemoryLimit:    getValue(mlps, i),
			CPUOnlineCount: getValue(cops, i),
			CPUUsage:       getValue(cups, i),
			CPUPercentage:  getValue(cpps, i),
		}
	}
	return resp, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func getValue(dps []*tstorage.DataPoint, index int) float64 {
	if len(dps) < index {
		return 0
	}

	return dps[index].Value
}
