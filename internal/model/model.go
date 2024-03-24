package model

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/eldius/docker-profiler/internal/helper"
	"log"
	"time"
)

const (
	kiloBytes = 1024
	megaBytes = 1048576
	gigaBytes = 1073741824
)

//type ContainerStats struct {
//	Read      time.Time `json:"read"`
//	Preread   time.Time `json:"preread"`
//	PidsStats struct {
//		Current int   `json:"current"`
//		Limit   int64 `json:"limit"`
//	} `json:"pids_stats"`
//	BlkioStats struct {
//		IoServiceBytesRecursive interface{} `json:"io_service_bytes_recursive"`
//		IoServicedRecursive     interface{} `json:"io_serviced_recursive"`
//		IoQueueRecursive        interface{} `json:"io_queue_recursive"`
//		IoServiceTimeRecursive  interface{} `json:"io_service_time_recursive"`
//		IoWaitTimeRecursive     interface{} `json:"io_wait_time_recursive"`
//		IoMergedRecursive       interface{} `json:"io_merged_recursive"`
//		IoTimeRecursive         interface{} `json:"io_time_recursive"`
//		SectorsRecursive        interface{} `json:"sectors_recursive"`
//	} `json:"blkio_stats"`
//	NumProcs     int `json:"num_procs"`
//	StorageStats struct {
//	} `json:"storage_stats"`
//	CPUStats struct {
//		CPUUsage struct {
//			TotalUsage        int `json:"total_usage"`
//			UsageInKernelmode int `json:"usage_in_kernelmode"`
//			UsageInUsermode   int `json:"usage_in_usermode"`
//		} `json:"cpu_usage"`
//		SystemCPUUsage int64 `json:"system_cpu_usage"`
//		OnlineCpus     int   `json:"online_cpus"`
//		ThrottlingData struct {
//			Periods          int `json:"periods"`
//			ThrottledPeriods int `json:"throttled_periods"`
//			ThrottledTime    int `json:"throttled_time"`
//		} `json:"throttling_data"`
//	} `json:"cpu_stats"`
//	PrecpuStats struct {
//		CPUUsage struct {
//			TotalUsage        int `json:"total_usage"`
//			UsageInKernelmode int `json:"usage_in_kernelmode"`
//			UsageInUsermode   int `json:"usage_in_usermode"`
//		} `json:"cpu_usage"`
//		SystemCPUUsage int64 `json:"system_cpu_usage"`
//		OnlineCpus     int   `json:"online_cpus"`
//		ThrottlingData struct {
//			Periods          int `json:"periods"`
//			ThrottledPeriods int `json:"throttled_periods"`
//			ThrottledTime    int `json:"throttled_time"`
//		} `json:"throttling_data"`
//	} `json:"precpu_stats"`
//	MemoryStats struct {
//		Usage int64 `json:"usage"`
//		Stats struct {
//			ActiveAnon            int `json:"active_anon"`
//			ActiveFile            int `json:"active_file"`
//			Anon                  int `json:"anon"`
//			AnonThp               int `json:"anon_thp"`
//			File                  int `json:"file"`
//			FileDirty             int `json:"file_dirty"`
//			FileMapped            int `json:"file_mapped"`
//			FileWriteback         int `json:"file_writeback"`
//			InactiveAnon          int `json:"inactive_anon"`
//			InactiveFile          int `json:"inactive_file"`
//			KernelStack           int `json:"kernel_stack"`
//			Pgactivate            int `json:"pgactivate"`
//			Pgdeactivate          int `json:"pgdeactivate"`
//			Pgfault               int `json:"pgfault"`
//			Pglazyfree            int `json:"pglazyfree"`
//			Pglazyfreed           int `json:"pglazyfreed"`
//			Pgmajfault            int `json:"pgmajfault"`
//			Pgrefill              int `json:"pgrefill"`
//			Pgscan                int `json:"pgscan"`
//			Pgsteal               int `json:"pgsteal"`
//			Shmem                 int `json:"shmem"`
//			Slab                  int `json:"slab"`
//			SlabReclaimable       int `json:"slab_reclaimable"`
//			SlabUnreclaimable     int `json:"slab_unreclaimable"`
//			Sock                  int `json:"sock"`
//			ThpCollapseAlloc      int `json:"thp_collapse_alloc"`
//			ThpFaultAlloc         int `json:"thp_fault_alloc"`
//			Unevictable           int `json:"unevictable"`
//			WorkingsetActivate    int `json:"workingset_activate"`
//			WorkingsetNodereclaim int `json:"workingset_nodereclaim"`
//			WorkingsetRefault     int `json:"workingset_refault"`
//		} `json:"stats"`
//		Limit int64 `json:"limit"`
//	} `json:"memory_stats"`
//	Name     string `json:"name"`
//	ID       string `json:"id"`
//	Networks struct {
//		Eth0 struct {
//			RxBytes   int `json:"rx_bytes"`
//			RxPackets int `json:"rx_packets"`
//			RxErrors  int `json:"rx_errors"`
//			RxDropped int `json:"rx_dropped"`
//			TxBytes   int `json:"tx_bytes"`
//			TxPackets int `json:"tx_packets"`
//			TxErrors  int `json:"tx_errors"`
//			TxDropped int `json:"tx_dropped"`
//		} `json:"eth0"`
//	} `json:"networks"`
//}

type ContainerStats struct {
	types.StatsJSON
}

func (s ContainerStats) MemoryUsageStr() string {
	return helper.FormatMemory(s.MemoryStats.Usage)
}

func (s ContainerStats) MemoryLimitStr() string {
	return fmt.Sprintf("%02.2f", float64(s.MemoryStats.Limit)/float64(1024*1024))
}

func (s *ContainerStats) CPUUsagePercentage() float64 {
	cpuPercent := 0.0
	//numCPUs := len(s.PercpuUsage)
	numCPUs := s.CPUStats.OnlineCPUs

	//cpuDelta := float64(s.CPUUSage) - float64(s.PreCPUUSage)
	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)

	//systemDelta := float64(s.SystemCPUUsage) - float64(s.PreSystemCPUUsage)
	systemDelta := float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)

	log.Printf("numCPUs: %d / cpuDelta: %01.2f / systemDelta: %01.2f\n", numCPUs, cpuDelta, systemDelta)
	log.Printf("numCPUs: %d / cpuDelta: %01.2f / systemDelta: %01.2f\n", numCPUs, cpuDelta, systemDelta)

	if cpuDelta > 0.0 && systemDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(numCPUs) * 100.0
	}
	return cpuPercent
}

type MetricsDatapoint struct {
	Timestamp      time.Time
	MemoryUsage    float64
	MemoryLimit    float64
	CPUOnlineCount float64
	CPUUsage       float64
	CPUPercentage  float64
}

func (m MetricsDatapoint) MemoryUsageStr() string {
	return helper.FormatMemory(uint64(m.MemoryUsage))
}

func (m MetricsDatapoint) MemoryLimitStr() string {
	return helper.FormatMemory(uint64(m.MemoryLimit))
}
