package docker

import "time"

type ContainerStats struct {
	Read      time.Time `json:"read"`
	Preread   time.Time `json:"preread"`
	PidsStats struct {
		Current int   `json:"current"`
		Limit   int64 `json:"limit"`
	} `json:"pids_stats"`
	BlkioStats struct {
		IoServiceBytesRecursive interface{} `json:"io_service_bytes_recursive"`
		IoServicedRecursive     interface{} `json:"io_serviced_recursive"`
		IoQueueRecursive        interface{} `json:"io_queue_recursive"`
		IoServiceTimeRecursive  interface{} `json:"io_service_time_recursive"`
		IoWaitTimeRecursive     interface{} `json:"io_wait_time_recursive"`
		IoMergedRecursive       interface{} `json:"io_merged_recursive"`
		IoTimeRecursive         interface{} `json:"io_time_recursive"`
		SectorsRecursive        interface{} `json:"sectors_recursive"`
	} `json:"blkio_stats"`
	NumProcs     int `json:"num_procs"`
	StorageStats struct {
	} `json:"storage_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int `json:"total_usage"`
			UsageInKernelmode int `json:"usage_in_kernelmode"`
			UsageInUsermode   int `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int `json:"total_usage"`
			UsageInKernelmode int `json:"usage_in_kernelmode"`
			UsageInUsermode   int `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage int `json:"usage"`
		Stats struct {
			ActiveAnon            int `json:"active_anon"`
			ActiveFile            int `json:"active_file"`
			Anon                  int `json:"anon"`
			AnonThp               int `json:"anon_thp"`
			File                  int `json:"file"`
			FileDirty             int `json:"file_dirty"`
			FileMapped            int `json:"file_mapped"`
			FileWriteback         int `json:"file_writeback"`
			InactiveAnon          int `json:"inactive_anon"`
			InactiveFile          int `json:"inactive_file"`
			KernelStack           int `json:"kernel_stack"`
			Pgactivate            int `json:"pgactivate"`
			Pgdeactivate          int `json:"pgdeactivate"`
			Pgfault               int `json:"pgfault"`
			Pglazyfree            int `json:"pglazyfree"`
			Pglazyfreed           int `json:"pglazyfreed"`
			Pgmajfault            int `json:"pgmajfault"`
			Pgrefill              int `json:"pgrefill"`
			Pgscan                int `json:"pgscan"`
			Pgsteal               int `json:"pgsteal"`
			Shmem                 int `json:"shmem"`
			Slab                  int `json:"slab"`
			SlabReclaimable       int `json:"slab_reclaimable"`
			SlabUnreclaimable     int `json:"slab_unreclaimable"`
			Sock                  int `json:"sock"`
			ThpCollapseAlloc      int `json:"thp_collapse_alloc"`
			ThpFaultAlloc         int `json:"thp_fault_alloc"`
			Unevictable           int `json:"unevictable"`
			WorkingsetActivate    int `json:"workingset_activate"`
			WorkingsetNodereclaim int `json:"workingset_nodereclaim"`
			WorkingsetRefault     int `json:"workingset_refault"`
		} `json:"stats"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	Networks struct {
		Eth0 struct {
			RxBytes   int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			RxErrors  int `json:"rx_errors"`
			RxDropped int `json:"rx_dropped"`
			TxBytes   int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
			TxErrors  int `json:"tx_errors"`
			TxDropped int `json:"tx_dropped"`
		} `json:"eth0"`
	} `json:"networks"`
}
