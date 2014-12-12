package client

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

const testAddress = "iot.eclipse.org:1883"

var errTest = errors.New("testError")

type errConn struct {
	net.TCPConn
}

func (c *errConn) Close() error {
	return errTest
}

type errPacket struct{}

func (p *errPacket) WriteTo(_ io.Writer) (int64, error) {
	return 0, errTest
}

func TestClient_Receive(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cli.sendCONNECT(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if _, _, err := cli.Receive(); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cli.Disconnect(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_Receive_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if _, _, err := cli.Receive(); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_Receive_errFirstReadByte(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cli.conn.R = bufio.NewReader(bytes.NewReader(nil))

	if _, _, err := cli.Receive(); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestClient_Receive_errSecondReadByte(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cli.conn.R = bufio.NewReader(bytes.NewReader([]byte{0x00}))

	if _, _, err := cli.Receive(); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestClient_Receive_errReadFull(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cli.conn.R = bufio.NewReader(bytes.NewReader([]byte{packet.TypeCONNACK << 4, 0x80, 0x01}))

	if _, _, err := cli.Receive(); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestClient_Receive_errNewCONNACKFromBytes(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cli.conn.R = bufio.NewReader(bytes.NewReader([]byte{packet.TypeCONNACK << 4, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))

	if _, _, err := cli.Receive(); err != packet.ErrCONNACKInvalidVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", packet.ErrCONNACKInvalidVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, packet.ErrCONNACKInvalidVariableHeaderLen)
		}
	}
}

func TestClient_establish_errAlreadyConnected(t *testing.T) {
	cli := &Client{
		conn: &mqtt.Connection{},
	}

	if err := cli.establish("tcp", ""); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_establish_errNewConnection(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", ""); err == nil {
		t.Error("err => nil, want => not nil")
	}
}
func TestClient_establish(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Error("err => %q, want => nil", err)
	}
}

func TestClient_close_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.close(); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_close_closeErr(t *testing.T) {
	cli := &Client{
		conn: &mqtt.Connection{
			Conn: &errConn{},
		},
	}

	if err := cli.close(); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestClient_close_cleanSession(t *testing.T) {
	cli := &Client{
		sess: NewSession(nil),
	}

	err := cli.establish("tcp", testAddress)

	if err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cli.Disconnect(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_sendCONNECT_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.sendCONNECT(nil); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_sendCONNECT_optsNil(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Error("err => %q, want => nil", err)
		return
	}

	if err := cli.sendCONNECT(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_sendCONNECT_reuseSession(t *testing.T) {
	var cleanSession bool

	cli := &Client{
		sess: NewSession(&SessionOptions{
			CleanSession: &cleanSession,
		}),
	}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Error("err => %q, want => nil", err)
		return
	}

	err := cli.sendCONNECT(&packet.CONNECTOptions{
		CleanSession: &cleanSession,
	})
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_sendCONNECT_newCONNECTErr(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Error("err => %q, want => nil", err)
		return

	}

	err := cli.sendCONNECT(&packet.CONNECTOptions{
		WillTopic: "willTopic",
	})
	if err != packet.ErrCONNECTWillTopicMessageEmpty {
		if err == nil {
			t.Errorf("err => nil, want => %q", packet.ErrCONNECTWillTopicMessageEmpty)
		} else {
			t.Errorf("err => %q, want => %q", err, packet.ErrCONNECTWillTopicMessageEmpty)
		}
	}
}

func TestClient_sendDISCONNECT_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.sendDISCONNECT(); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_sendDISCONNECT(t *testing.T) {
	cli := &Client{}

	if err := cli.establish("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cli.sendDISCONNECT(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_send_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.send(packet.NewDISCONNECT()); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_send_errWriteTo(t *testing.T) {
	cli := &Client{
		conn: &mqtt.Connection{},
	}

	if err := cli.send(&errPacket{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestNew(t *testing.T) {
	cli := New(nil)

	if cli == nil {
		t.Error("cli => nil, want => not nil")
	}
}
