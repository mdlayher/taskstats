//+build !linux

package taskstats

import (
	"fmt"
	"runtime"
)

var (
	// errUnimplemented is returned by all functions on platforms that
	// cannot make use of taskstats.
	errUnimplemented = fmt.Errorf("taskstats not implemented on %s/%s",
		runtime.GOOS, runtime.GOARCH)
)

var _ osClient = &client{}

// A client is an unimplemented taskstats client.
type client struct{}

// newClient always returns an error.
func newClient() (*client, error) {
	return nil, errUnimplemented
}

// Close implements osClient.
func (c *client) Close() error {
	return errUnimplemented
}

// CGroupStats implements osClient.
func (c *client) CGroupStats(path string) (*CGroupStats, error) {
	return nil, errUnimplemented
}

// PID implements osClient.
func (c *client) PID(pid int) (*Stats, error) {
	return nil, errUnimplemented
}

// TGID implements osClient.
func (c *client) TGID(tgid int) (*Stats, error) {
	return nil, errUnimplemented
}
