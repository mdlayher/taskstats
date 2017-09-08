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

			// TODO(mdlayher): parse raw unix.Taskstats structure into nicer structure.
			stats := Stats(*(*unix.Taskstats)(unsafe.Pointer(&na.Data[0])))
			return &stats, nil
		}
	}

	// No taskstats response found.
	return nil, os.ErrNotExist
}
