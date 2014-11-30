package client

import (
	"errors"
	"io"
	"testing"
	"time"
)

const testAddress = "iot.eclipse.org:1883"

type packetErr struct{}

func (w *packetErr) WriteTo(_ io.Writer) (int64, error) {
	return 0, errTest
}

type readerErr struct{}

func (r *readerErr) Read(_ []byte) (int, error) {
	return 0, errTest
}

var errTest = errors.New("test error")

func TestClient_Connect_addressInvalid(t *testing.T) {
	if err := New().Connect(&ConnectOptions{Address: "test"}, nil); err == nil {
		t.Errorf("err => nil, want => %q", err)
	}
}

func TestClient_Connect_errAlreadyConnected(t *testing.T) {
	cli := New()

	if err := cli.Connect(&ConnectOptions{Address: testAddress}, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Connect(&ConnectOptions{Address: testAddress}, nil); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Connect(t *testing.T) {
	if err := New().Connect(&ConnectOptions{Address: testAddress}, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_Connect_sendErr(t *testing.T) {
	cli := New()

	if err := cli.Connect(&ConnectOptions{Address: testAddress}, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cli.sendc <- &packetErr{}

	select {
	case err := <-cli.Errc:
		if err != errTest {
			if err == nil {
				t.Errorf("err => nil, want => %q", errTest)
			} else {
				t.Errorf("err => %q, want => %q", err, errTest)
			}
		}
	case <-time.After(10 * time.Second):
		t.Errorf("err => nil, want => %q", errTest)
	}
}
