package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/yosssi/gmq/common/packet"
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

func TestClient_Connect_addressEmpty(t *testing.T) {
	if err := New().Connect("", nil); err == nil {
		t.Errorf("err => nil, want => %q", err)
	}
}

func TestClient_Connect_errAlreadyConnected(t *testing.T) {
	cli := New()

	if err := cli.Connect(testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Connect(testAddress, nil); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Connect(t *testing.T) {
	if err := New().Connect(testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestClient_Connect_sendErr(t *testing.T) {
	cli := New()

	if err := cli.Connect(testAddress, nil); err != nil {
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

func TestClient_Connect_receive(t *testing.T) {
	cli := New()

	if err := cli.Connect(testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cli.conn.Write([]byte{0})

	select {
	case <-cli.Errc:
	case <-time.After(10 * time.Second):
		t.Error("err => nil, want => not nil")
	}

}

func TestClient_send_err(t *testing.T) {
	cli := New()

	w := bufio.NewWriter(ioutil.Discard)

	if err := cli.send(w, &packetErr{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestClient_receive_err(t *testing.T) {
	cli := New()

	r := bufio.NewReader(&readerErr{})

	if err := cli.receive(r); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestClient_receive_errInvalidHeader(t *testing.T) {
	cli := New()

	b := make([]byte, 268435457)
	b[0] = 32
	b[1] = 128
	b[2] = 0xF0
	b[3] = 0xF0
	b[4] = 0xF0

	r := bufio.NewReader(bytes.NewReader(b))

	if err := cli.receive(r); err != packet.ErrInvalidCONNACKVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", packet.ErrInvalidCONNACKVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, packet.ErrInvalidCONNACKVariableHeaderLen)
		}
	}
}

func TestClient_receive_errReadRemaining(t *testing.T) {
	cli := New()

	b := make([]byte, 5)
	b[0] = 32
	b[1] = 128
	b[2] = 0xF0
	b[3] = 0xF0
	b[4] = 0xF0

	r := bufio.NewReader(bytes.NewReader(b))

	if err := cli.receive(r); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestClient_receive_errReadFull(t *testing.T) {
	cli := New()

	b := make([]byte, 6)
	b[0] = 32
	b[1] = 128
	b[2] = 0xF0
	b[3] = 0xF0
	b[4] = 0
	b[5] = 0xF0

	r := bufio.NewReader(bytes.NewReader(b))

	err := cli.receive(r)
	fmt.Println(err)
}

func TestNew(t *testing.T) {
	if cli := New(); cli == nil {
		t.Error("cli should not be nil")
	}
}
