//+build linux

package taskstats_test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/taskstats"
	"golang.org/x/sys/unix"
)

func TestLinuxClientIntegration(t *testing.T) {
	c, err := taskstats.New()
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}
	defer c.Close()

	stats, err := c.Self()
	if err != nil {
		if os.IsPermission(err) {
			t.Skipf("taskstats requires elevated permission: %v", err)
		}

		t.Fatalf("failed to retrieve self stats: %v", err)
	}

	if diff := cmp.Diff(unix.TASKSTATS_VERSION, int(stats.Version)); diff != "" {
		t.Fatalf("unexpected taskstats version (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(os.Getpid(), int(stats.Ac_pid)); diff != "" {
		t.Fatalf("unexpected PID (-want +got):\n%s", diff)
	}

	// TODO(mdlayher): verify more fields?
}
