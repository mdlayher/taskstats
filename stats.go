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

// TODO(mdlayher): export cleaned up Stats type.
type stats struct {
	BeginTime       time.Time
	ElapsedTime     time.Duration
	UserCPUTime     time.Duration
	SystemCPUTime   time.Duration
	MinorPageFaults uint64
	MajorPageFaults uint64

	CPUCount uint64
	CPUDelay time.Duration
}
