//+build linux

package taskstats

import (
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/genetlink/genltest"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func TestLinuxClientPIDBadMessages(t *testing.T) {
	tests := []struct {
		name string
		msgs []genetlink.Message
	}{
		{
			name: "no messages",
			msgs: []genetlink.Message{},
		},
		{
			name: "two messages",
			msgs: []genetlink.Message{{}, {}},
		},
		{
			name: "incorrect taskstats size",
			msgs: []genetlink.Message{{
				Data: mustMarshalAttributes([]netlink.Attribute{{
					Type: unix.TASKSTATS_TYPE_AGGR_PID,
					Data: mustMarshalAttributes([]netlink.Attribute{{
						Type: unix.TASKSTATS_TYPE_STATS,
						Data: []byte{0xff},
					}}),
				}}),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient(t, func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			})
			defer c.Close()

			_, err := c.PID(1)
			if err == nil {
				t.Fatal("an error was expected, but none occurred")
			}
		})
	}
}

func TestLinuxClientPIDIsNotExist(t *testing.T) {
	tests := []struct {
		name string
		msg  genetlink.Message
	}{
		{
			name: "no attributes",
			msg:  genetlink.Message{},
		},
		{
			name: "no aggr+pid",
			msg: genetlink.Message{
				Data: mustMarshalAttributes([]netlink.Attribute{{
					Type: unix.TASKSTATS_TYPE_NULL,
				}}),
			},
		},
		{
			name: "no stats",
			msg: genetlink.Message{
				Data: mustMarshalAttributes([]netlink.Attribute{{
					Type: unix.TASKSTATS_TYPE_AGGR_PID,
					Data: mustMarshalAttributes([]netlink.Attribute{{
						Type: unix.TASKSTATS_TYPE_NULL,
					}}),
				}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient(t, func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return []genetlink.Message{tt.msg}, nil
			})
			defer c.Close()

			_, err := c.PID(1)
			if !os.IsNotExist(err) {
				t.Fatalf("expected is not exist, but got: %v", err)
			}
		})
	}
}

func TestLinuxClientPIDOK(t *testing.T) {
	pid := os.Getpid()

	stats := unix.Taskstats{
		Version: unix.TASKSTATS_VERSION,
		Ac_pid:  uint32(pid),
	}

	fn := func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
		// Cast unix.Taskstats structure into a byte array with the correct size.
		b := *(*[unsafe.Sizeof(stats)]byte)(unsafe.Pointer(&stats))

		return []genetlink.Message{{
			Data: mustMarshalAttributes([]netlink.Attribute{{
				Type: unix.TASKSTATS_TYPE_AGGR_PID,
				Data: mustMarshalAttributes([]netlink.Attribute{{
					Type: unix.TASKSTATS_TYPE_STATS,
					Data: b[:],
				}}),
			}}),
		}}, nil
	}

	c := testClient(t, checkRequest(unix.TASKSTATS_CMD_GET, netlink.HeaderFlagsRequest, fn))
	defer c.Close()

	newStats, err := c.PID(pid)
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}

	tstats := Stats(stats)

	if diff := cmp.Diff(&tstats, newStats); diff != "" {
		t.Fatalf("unexpected taskstats structure (-want +got):\n%s", diff)
	}
}

func checkRequest(command uint8, flags netlink.HeaderFlags, fn genltest.Func) genltest.Func {
	return func(greq genetlink.Message, nreq netlink.Message) ([]genetlink.Message, error) {
		if want, got := command, greq.Header.Command; command != 0 && want != got {
			return nil, fmt.Errorf("unexpected generic netlink header command: %d, want: %d", got, want)
		}

		if want, got := flags, nreq.Header.Flags; flags != 0 && want != got {
			return nil, fmt.Errorf("unexpected generic netlink header command: %s, want: %s", got, want)
		}

		return fn(greq, nreq)
	}
}

func testClient(t *testing.T, fn genltest.Func) *client {
	family := genetlink.Family{
		ID:      20,
		Version: unix.TASKSTATS_GENL_VERSION,
		Name:    unix.TASKSTATS_GENL_NAME,
	}

	conn := genltest.Dial(genltest.ServeFamily(family, fn))

	c, err := initClient(conn)
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}

	return c
}

func mustMarshalAttributes(attrs []netlink.Attribute) []byte {
	b, err := netlink.MarshalAttributes(attrs)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal attributes: %v", err))
	}

	return b
}
