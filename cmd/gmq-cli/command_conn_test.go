package main

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

const testAddress = "iot.eclipse.org:1883"

var errTest = errors.New("test")

type packetErr struct{}

func (p packetErr) WriteTo(w io.Writer) (int64, error) {
	return 0, errTest
}

func (p packetErr) Type() (byte, error) {
	return 0x00, errTest
}

func Test_commandConn_run_err(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cmd.run(); err == nil {
		t.Error("err => nil, want => not nil")
	}
}

func Test_commandConn_run(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.network = "tcp"
	cmd.address = testAddress

	if err := cmd.run(); err != nil {
		t.Error("err => %q, want => nil", err)
	}

	if err := disconnect(cmd.ctx); err != nil {
		t.Error("err => %q, want => nil", err)
	}
}

func Test_commandConn_waitCONNACK_connack(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.waitCONNACK()

	cmd.ctx.connack <- struct{}{}

	cmd.ctx.wg.Wait()
}

func Test_commandConn_waitCONNACK_timeout(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.connackTimeout = 1

	cmd.ctx.wg.Add(1)
	go cmd.waitCONNACK()

	cmd.ctx.wg.Wait()
}

func Test_commandConn_waitCONNACK_timeout_disconnDefault(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.connackTimeout = 1

	cmd.ctx.disconn <- struct{}{}

	cmd.ctx.wg.Add(1)
	go cmd.waitCONNACK()

	cmd.ctx.wg.Wait()
}

func Test_commandConn_waitCONNACK_connackEnd(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.waitCONNACK()

	cmd.ctx.connackEnd <- struct{}{}

	cmd.ctx.wg.Wait()
}

func Test_commandConn_receive_ReceiveErr_disconnecting(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.disconnecting = true

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	cmd.ctx.wg.Wait()
}

func Test_commandConn_receive_ReceiveErr(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	cmd.ctx.wg.Wait()
}

func Test_commandConn_receive_ReceiveErr_default(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.disconn <- struct{}{}

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	cmd.ctx.wg.Wait()
}

func Test_commandConn_receive_handleErr(t *testing.T) {
	defer func(handleBak func(cmd *commandConn, p packet.Packet) error) {
		handle = handleBak
	}(handle)

	handle = func(_ *commandConn, _ packet.Packet) error {
		return errTest
	}

	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	ctx.cli = client.New(nil)

	if err := ctx.cli.Connect("tcp", testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	time.Sleep(time.Second)

	if err := disconnect(cmd.ctx); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Wait()
}

func Test_commandConn_receive(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	ctx.cli = client.New(nil)

	if err := ctx.cli.Connect("tcp", testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	time.Sleep(time.Second)

	if err := disconnect(cmd.ctx); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Wait()
}

func Test_commandConn_handle_err(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cmd.handle(packetErr{}); err != errTest {
		errorfErr(t, err, errTest)
	}
}
