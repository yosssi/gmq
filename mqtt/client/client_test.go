package client

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/yosssi/gmq/mqtt"
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

func TestClient_handlePINGRESP_ErrInvalidPINGRESP(t *testing.T) {
	cli := New(&Options{
		ErrHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	if err := cli.handlePINGRESP(); err != ErrInvalidPINGRESP {
		invalidError(t, err, ErrInvalidPINGRESP)
	}
}

func TestClient_handlePINGRESP_default(t *testing.T) {
	cli := New(&Options{
		ErrHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}))

	if err := cli.handlePINGRESP(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePINGRESP(t *testing.T) {
	cli := New(&Options{
		ErrHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}, 1))

	if err := cli.handlePINGRESP(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleErrorAndDisconn_connNil(t *testing.T) {
	cli := New(nil)

	cli.handleErrorAndDisconn(errTest)
}

func TestClient_handleErrorAndDisconn(t *testing.T) {
	cli := New(&Options{
		ErrHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.handleErrorAndDisconn(errTest)

	time.Sleep(1 * time.Second)
}

func TestClient_sendPackets_sendErr(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.conn.Close(); err != nil {
		nilErrorExpected(t, err)
	}

	cli.conn.send <- &packetErr{}

	cli.conn.wg.Wait()
}

func TestClient_sendPackets_keepAliveErr(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:   "tcp",
		Address:   testAddress,
		ClientID:  []byte("clientID"),
		KeepAlive: 2,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.conn.Close(); err != nil {
		nilErrorExpected(t, err)
	}

	cli.conn.wg.Wait()
}

func TestClient_sendPackets_sendEnd(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:   "tcp",
		Address:   testAddress,
		ClientID:  []byte("clientID"),
		KeepAlive: 2,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	time.Sleep(3 * time.Second)

	cli.muConn.Lock()
	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}))
	cli.muConn.Unlock()

	cli.conn.sendEnd <- struct{}{}

	cli.conn.wg.Wait()
}

func TestClient_newPUBLISHPacket_generatePacketID(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	id := minPacketID

	for {
		cli.sess.sendingPackets[id] = nil

		if id == maxPacketID {
			break
		}

		id++
	}

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: mqtt.QoS1,
	})

	if err != ErrPacketIDExhaused {
		invalidError(t, err, ErrPacketIDExhaused)
	}
}

func TestClient_newPUBLISHPacket_ErrInvalidQoS(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: 0x03,
	})

	if err != packet.ErrInvalidQoS {
		invalidError(t, err, packet.ErrInvalidQoS)
	}
}

func TestClient_newPUBLISHPacket(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: mqtt.QoS1,
	})

	if err != nil {
		nilErrorExpected(t, err)
	}
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
