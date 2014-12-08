package client

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/yosssi/gmq/mqtt/packet"
)

func TestReceive(t *testing.T) {
	cli := &Client{}

	err := cli.Connect(
		&ConnectOptions{
			Address: testAddress,
		},
	)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	err = cli.SendCONNECT(nil)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	_, _, err = Receive(cli.Conn.R)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Disconnect(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestReceive_errFirstReadByte(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader(nil))

	if _, _, err := Receive(r); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestReceive_errSecondReadByte(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0x00}))

	if _, _, err := Receive(r); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestReceive_errReadFull(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{packet.TypeCONNACK << 4, 0x80, 0x01}))

	if _, _, err := Receive(r); err != io.EOF {
		if err == nil {
			t.Errorf("err => nil, want => %q", io.EOF)
		} else {
			t.Errorf("err => %q, want => %q", err, io.EOF)
		}
	}
}

func TestReceive_errNewCONNACKFromBytes(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{packet.TypeCONNACK << 4, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))

	if _, _, err := Receive(r); err != packet.ErrCONNACKInvalidVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", packet.ErrCONNACKInvalidVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, packet.ErrCONNACKInvalidVariableHeaderLen)
		}
	}
}
