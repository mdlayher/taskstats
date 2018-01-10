package taskstats

import (
	"time"

	"golang.org/x/sys/unix"
)

// TODO(mdlayher+andrestc): remove this Stats type when stats is exported.

// Stats is the structure returned by Linux's taskstats interface.
type Stats unix.Taskstats

// parseCGroupStats parses a raw cgroupstats structure into a cleaner form.
func parseCGroupStats(cs unix.CGroupStats) (*CGroupStats, error) {
	// This conversion isn't really necessary for this type, but it allows us
	// to export a structure that isn't defined in a platform-specific way.
	stats := &CGroupStats{
		Sleeping:        cs.Sleeping,
		Running:         cs.Running,
		Stopped:         cs.Stopped,
		Uninterruptible: cs.Uninterruptible,
		IOWait:          cs.Io_wait,
	}

	return stats, nil
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
