//go:build linux
// +build linux

package taskstats_test

import (
	"os"
	"testing"

	"github.com/mdlayher/taskstats"
)

func TestLinuxClientIntegration(t *testing.T) {
	c, err := taskstats.New()
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}
	defer c.Close()

	t.Run("self", func(t *testing.T) {
		testSelfStats(t, c)
	})

	t.Run("cgroup", func(t *testing.T) {
		testCGroupStats(t, c)
	})
}

func testSelfStats(t *testing.T, c *taskstats.Client) {
	stats, err := c.Self()
	if err != nil {
		if os.IsPermission(err) {
			t.Skipf("taskstats requires elevated permission: %v", err)
		}

		t.Fatalf("failed to retrieve self stats: %v", err)
	}

	if stats.BeginTime.IsZero() {
		t.Fatalf("unexpected zero begin time")
	}

	// TODO(mdlayher): verify more fields?
}

func testCGroupStats(t *testing.T, c *taskstats.Client) {
	// TODO(mdlayher): try to verify these in some meaningful way, but for now,
	// no error means the structure is valid, which works.
	_, err := c.CGroupStats("/sys/fs/cgroup/cpu")
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		t.Skipf("did not find cgroup CPU stats: %v", err)
	}

	t.Fatalf("failed to retrieve cgroup stats: %v", err)
}
