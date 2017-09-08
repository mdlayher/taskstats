package taskstats

import (
	"time"

	"golang.org/x/sys/unix"
)

// Stats is the structure returned by Linux's taskstats interface.
type Stats unix.Taskstats

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

// parseStats parses a raw taskstats structure into a cleaner form.
func parseStats(ts unix.Taskstats) (*stats, error) {
	stats := &stats{
		BeginTime:       time.Unix(int64(ts.Ac_btime), 0),
		ElapsedTime:     microseconds(ts.Ac_etime),
		UserCPUTime:     microseconds(ts.Ac_utime),
		SystemCPUTime:   microseconds(ts.Ac_stime),
		MinorPageFaults: ts.Ac_minflt,
		MajorPageFaults: ts.Ac_majflt,

		CPUCount: ts.Cpu_count,
		CPUDelay: nanoseconds(ts.Cpu_delay_total),
	}

	return stats, nil
}

// nanoseconds converts a raw number of nanoseconds into a time.Duration.
func nanoseconds(t uint64) time.Duration {
	return time.Duration(t) * time.Nanosecond
}

// microseconds converts a raw number of microseconds into a time.Duration.
func microseconds(t uint64) time.Duration {
	return time.Duration(t) * time.Microsecond
}
