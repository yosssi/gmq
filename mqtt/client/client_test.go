package client

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/yosssi/gmq/mqtt/packet"
)

var errTest = errors.New("test error")

type packetErr struct{}

func (p *packetErr) WriteTo(w io.Writer) (int64, error) {
	return 0, errTest
}

func (p *packetErr) Type() (byte, error) {
	return 0x00, errTest
}

func TestClient_validatePacketID_notExist(t *testing.T) {
	cli := New(nil)

	if err := cli.validatePacketID(map[uint16]packet.Packet{}, 1, packet.TypePUBLISH); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestClient_validatePacketID_ptypeErr(t *testing.T) {
	cli := New(nil)

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: &packetErr{}}, 1, packet.TypePUBLISH); err != errTest {
		invalidError(t, err, errTest)
	}
}

func TestClient_validatePacketID_ErrInvalidPacketID(t *testing.T) {
	cli := New(nil)

	p, err := packet.NewPUBACK(&packet.PUBACKOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: p}, 1, packet.TypePUBLISH); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestClient_validatePacketID(t *testing.T) {
	cli := New(nil)

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: p}, 1, packet.TypePUBLISH); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleMessage_continue(t *testing.T) {
	cli := New(nil)

	cli.conn = &connection{}

	cli.conn.ackedSubs = map[string]MessageHandler{
		"test": nil,
	}

	cli.handleMessage([]byte("test"), nil)
}

func TestClient_handleMessage(t *testing.T) {
	cli := New(nil)

	cli.conn = &connection{}

	cli.conn.ackedSubs = map[string]MessageHandler{
		"test": func(_, _ []byte) {},
	}

	cli.handleMessage([]byte("test"), nil)
}

func TestNew_optsNil(t *testing.T) {
	cli := New(nil)

	cli.disconnc <- struct{}{}

	time.Sleep(500 * time.Millisecond)

	cli.disconnEndc <- struct{}{}

	cli.wg.Wait()

}

func TestNew(t *testing.T) {
	cli := New(&Options{
		ErrHandler: func(_ error) {},
	})

	cli.disconnc <- struct{}{}

	time.Sleep(500 * time.Millisecond)

	cli.disconnEndc <- struct{}{}

	cli.wg.Wait()
}

func Test_match(t *testing.T) {
	testCases := []struct {
		in struct {
			topicName   string
			topicFilter string
		}
		out bool
	}{
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/score/wimbledon",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/player1/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis",
				topicFilter: "sport/tennis/+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis2/",
				topicFilter: "+/tennis/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport//player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/test/player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/player1",
				topicFilter: "sport/+/player1",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "+/monitor/Clients",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients/test",
				topicFilter: "$SYS/monitor/+",
			},
			out: false,
		},
	}

	for _, tc := range testCases {
		if got := match(tc.in.topicName, tc.in.topicFilter); got != tc.out {
			t.Errorf("got => %t, want => %t", got, tc.out)
		}
	}
}

func invalidError(t *testing.T, err, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
