// Package taskstats provides access to Linux's taskstats interface, for sending
// per-task and per-process statistics from the kernel to userspace.
//
// For more information on taskstats, please see:
//   - https://www.kernel.org/doc/Documentation/accounting/taskstats.txt
//   - https://www.kernel.org/doc/Documentation/accounting/taskstats-struct.txt
//   - https://andrestc.com/post/linux-delay-accounting/
package taskstats

import (
	"io"
	"os"
)

// A Client provides access to Linux taskstats information.  Client operations
// require elevated privileges.
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

// Self is a convenience method for retrieving statistics about the current
// process.
func (c *Client) Self() (*Stats, error) {
	return c.c.PID(os.Getpid())
}

// PID retrieves statistics about a process, identified by its PID.
func (c *Client) PID(pid int) (*Stats, error) {
	return c.c.PID(pid)
}

// Close releases resources used by a Client.
func (c *Client) Close() error {
	return c.c.Close()
}

// An osClient is the operating system-specific implementation of Client.
type osClient interface {
	io.Closer
	PID(pid int) (*Stats, error)
}
