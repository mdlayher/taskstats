package taskstats

import "time"

// CGroupStats contains statistics for tasks of an individual cgroup.
type CGroupStats struct {
	Sleeping        uint64
	Running         uint64
	Stopped         uint64
	Uninterruptible uint64
	IOWait          uint64
}

// Stats contains statistics for an individual task.
type Stats struct {
	BeginTime           time.Time
	ElapsedTime         time.Duration
	UserCPUTime         time.Duration
	SystemCPUTime       time.Duration
	MinorPageFaults     uint64
	MajorPageFaults     uint64
	CPUDelayCount       uint64
	CPUDelay            time.Duration
	BlockIODelayCount   uint64
	BlockIODelay        time.Duration
	SwapInDelayCount    uint64
	SwapInDelay         time.Duration
	FreePagesDelayCount uint64
	FreePagesDelay      time.Duration
}
