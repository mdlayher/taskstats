// Package taskstats provides access to Linux's taskstats interface, for sending
// per-task, per-process, and cgroup statistics from the kernel to userspace.
//
// For more information on taskstats, please see:
//   - https://www.kernel.org/doc/Documentation/accounting/cgroupstats.txt
//   - https://www.kernel.org/doc/Documentation/accounting/taskstats.txt
//   - https://www.kernel.org/doc/Documentation/accounting/taskstats-struct.txt
//   - https://andrestc.com/post/linux-delay-accounting/
package taskstats

import (
	"io"
	"os"
)

// A Client provides access to Linux taskstats information.
//
// Some Client operations require elevated privileges.
type Client struct {
	c osClient
}

// New creates a new Client.
func New() (*Client, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

// CGroupStats retrieves cgroup statistics for the cgroup specified by path.
// Path should be a CPU cgroup path found in sysfs, such as:
//  - /sys/fs/cgroup/cpu
//  - /sys/fs/cgroup/cpu/docker
//  - /sys/fs/cgroup/cpu/docker/(hexadecimal identifier)
func (c *Client) CGroupStats(path string) (*CGroupStats, error) {
	return c.c.CGroupStats(path)
}

// Self is a convenience method for retrieving statistics about the current
// process.
func (c *Client) Self() (*Stats, error) {
	return c.c.TGID(os.Getpid())
}

// PID retrieves statistics about a process, identified by its PID.
func (c *Client) PID(pid int) (*Stats, error) {
	return c.c.PID(pid)
}

// TGID retrieves statistics about a thread group, identified by its TGID.
func (c *Client) TGID(tgid int) (*Stats, error) {
	return c.c.TGID(tgid)
}

// Close releases resources used by a Client.
func (c *Client) Close() error {
	return c.c.Close()
}

// An osClient is the operating system-specific implementation of Client.
type osClient interface {
	io.Closer
	CGroupStats(path string) (*CGroupStats, error)
	PID(pid int) (*Stats, error)
	TGID(tgid int) (*Stats, error)
}
