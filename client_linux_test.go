//+build linux

package taskstats_test

import (
	"log"
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

	stats, err := c.Self()
	if err != nil {
		if os.IsPermission(err) {
			t.Skipf("taskstats requires elevated permission: %v", err)
		}

		t.Fatalf("failed to retrieve self stats: %v", err)
	}

	// TODO(mdlayher): verify fields
	log.Println(stats)
}
