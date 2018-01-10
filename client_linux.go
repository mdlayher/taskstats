//+build linux

package taskstats

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/sys/unix"
)

const (
	// sizeofTaskstats is the size of a unix.Taskstats structure.
	sizeofTaskstats   = int(unsafe.Sizeof(unix.Taskstats{}))
	sizeofCGroupStats = int(unsafe.Sizeof(unix.CGroupStats{}))
)

var _ osClient = &client{}

// A client is a Linux-specific taskstats client.
type client struct {
	c      *genetlink.Conn
	family genetlink.Family
}

// newClient opens a connection to the taskstats family using
// generic netlink.
func newClient() (*client, error) {
	c, err := genetlink.Dial(nil)
	if err != nil {
		return nil, err
	}

	return initClient(c)
}

// initClient is the internal client constructor used in some tests.
func initClient(c *genetlink.Conn) (*client, error) {
	f, err := c.GetFamily(unix.TASKSTATS_GENL_NAME)
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	return &client{
		c:      c,
		family: f,
	}, nil
}

// Close implements osClient.
func (c *client) Close() error {
	return c.c.Close()
}

// PID implements osClient.
func (c *client) PID(pid int) (*Stats, error) {
	// Query taskstats for information using a specific PID.
	attrb, err := netlink.MarshalAttributes([]netlink.Attribute{{
		Type: unix.TASKSTATS_CMD_ATTR_PID,
		Data: nlenc.Uint32Bytes(uint32(pid)),
	}})
	if err != nil {
		return nil, err
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: unix.TASKSTATS_CMD_GET,
			Version: unix.TASKSTATS_VERSION,
		},
		Data: attrb,
	}

	flags := netlink.HeaderFlagsRequest

	msgs, err := c.c.Execute(msg, c.family.ID, flags)
	if err != nil {
		return nil, err
	}

	if l := len(msgs); l != 1 {
		return nil, fmt.Errorf("unexpected number of taskstats messages: %d", l)
	}

	return parseMessage(msgs[0])
}

// CGroupStats implements osClient.
func (c *client) CGroupStats(path string) (*CGroupStats, error) {
	// Open cgroup path so its file descriptor can be passed to taskstats.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Query taskstats for cgroup information using the file descriptor.
	attrb, err := netlink.MarshalAttributes([]netlink.Attribute{{
		Type: unix.CGROUPSTATS_CMD_ATTR_FD,
		Data: nlenc.Uint32Bytes(uint32(f.Fd())),
	}})
	if err != nil {
		return nil, err
	}

	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: unix.CGROUPSTATS_CMD_GET,
			Version: unix.TASKSTATS_VERSION,
		},
		Data: attrb,
	}

	flags := netlink.HeaderFlagsRequest

	msgs, err := c.c.Execute(msg, c.family.ID, flags)
	if err != nil {
		return nil, err
	}

	if l := len(msgs); l != 1 {
		return nil, fmt.Errorf("unexpected number of cgroupstats messages: %d", l)
	}

	return parseCGroupMessage(msgs[0])
}

// parseCGroupMessage attempts to parse a CGroupStats structure from a generic netlink message.
func parseCGroupMessage(m genetlink.Message) (*CGroupStats, error) {
	attrs, err := netlink.UnmarshalAttributes(m.Data)
	if err != nil {
		return nil, err
	}

	for _, a := range attrs {
		// Only parse cgroupstats structure.
		if a.Type != unix.CGROUPSTATS_TYPE_CGROUP_STATS {
			continue
		}

		// Verify that the byte slice containing a unix.CGroupStats is the
		// size expected by this package, so we don't blindly cast the
		// byte slice into a structure of the wrong size.
		if want, got := sizeofCGroupStats, len(a.Data); want != got {
			return nil, fmt.Errorf("unexpected cgroupstats structure size, want %d, got %d", want, got)
		}

		cs := *(*unix.CGroupStats)(unsafe.Pointer(&a.Data[0]))
		return parseCGroupStats(cs)
	}

	// No taskstats response found.
	return nil, os.ErrNotExist
}

// parseMessage attempts to parse a Stats structure from a generic netlink message.
func parseMessage(m genetlink.Message) (*Stats, error) {
	attrs, err := netlink.UnmarshalAttributes(m.Data)
	if err != nil {
		return nil, err
	}

	for _, a := range attrs {
		// Only parse PID+stats structure.
		if a.Type != unix.TASKSTATS_TYPE_AGGR_PID {
			continue
		}

		nattrs, err := netlink.UnmarshalAttributes(a.Data)
		if err != nil {
			return nil, err
		}

		for _, na := range nattrs {
			// Only parse Stats element since caller would already have PID.
			if na.Type != unix.TASKSTATS_TYPE_STATS {
				continue
			}

			// Verify that the byte slice containing a unix.Taskstats is the
			// size expected by this package, so we don't blindly cast the
			// byte slice into a structure of the wrong size.
			if want, got := sizeofTaskstats, len(na.Data); want != got {
				return nil, fmt.Errorf("unexpected taskstats structure size, want %d, got %d", want, got)
			}

			return parseStats(*(*unix.Taskstats)(unsafe.Pointer(&na.Data[0])))
		}
	}

	// No taskstats response found.
	return nil, os.ErrNotExist
}
