package client

import (
	"errors"
	"io"
	"net"
	"testing"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

const testAddress = "iot.eclipse.org:1883"

var errTest = errors.New("testError")

type errWriter struct{}

func (w *errWriter) Write(_ []byte) (int, error) {
	return 0, errTest
}

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

func TestClient_Connect_errAlreadyConnected(t *testing.T) {
	cli := &Client{
		Conn: &mqtt.Connection{},
	}

	if err := cli.Connect(nil); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Connect_optsNil(t *testing.T) {
	cli := &Client{}

	if err := cli.Connect(nil); err == nil {
		t.Error("err => nil, want => not nil")
	}
}

func TestClient_Connect(t *testing.T) {
	cli := &Client{}

	err := cli.Connect(&ConnectOptions{
		Address: testAddress,
	})
	if err != nil {
		t.Error("err => %q, want => nil", err)
	}
}

func TestClient_Disconnect_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.Disconnect(); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_Disconnect_closeErr(t *testing.T) {
	cli := &Client{
		Conn: &mqtt.Connection{
			Conn: &errConn{},
		},
	}

	if err := cli.Disconnect(); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestClient_Disconnect_cleanSession(t *testing.T) {
	cli := &Client{
		Sess: NewSession(nil),
	}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Disconnect(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_SendCONNECT_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.SendCONNECT(nil); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_SendCONNECT_optsNil(t *testing.T) {
	cli := &Client{}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)
	if err != nil {
		t.Error("err => %q, want => nil", err)
	}

	if err := cli.SendCONNECT(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_SendCONNECT_reuseSession(t *testing.T) {
	var cleanSession bool

	cli := &Client{
		Sess: NewSession(&SessionOptions{
			CleanSession: &cleanSession,
		}),
	}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)

	if err != nil {
		t.Error("err => %q, want => nil", err)
	}

	err = cli.SendCONNECT(&packet.CONNECTOptions{
		CleanSession: &cleanSession,
	})
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_Connect_newCONNECTErr(t *testing.T) {
	cli := &Client{}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)
	if err != nil {
		t.Error("err => %q, want => nil", err)
	}

	err = cli.SendCONNECT(&packet.CONNECTOptions{
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

func TestClient_SendDISCONNECT_errNotYetConnected(t *testing.T) {
	cli := &Client{}

	if err := cli.SendDISCONNECT(); err != ErrNotYetConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrNotYetConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrNotYetConnected)
		}
	}
}

func TestClient_SendDISCONNECT(t *testing.T) {
	cli := &Client{}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)

	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.SendDISCONNECT(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_Send_err(t *testing.T) {
	cli := &Client{
		Conn: &mqtt.Connection{},
	}

	if err := cli.send(&errPacket{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}
