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

// Stats is not implemented on this platform.
type Stats struct{}

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

// PID implements osClient.
func (c *client) PID(pid int) (*Stats, error) {
	return nil, errUnimplemented
}
