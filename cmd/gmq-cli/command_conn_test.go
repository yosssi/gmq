package main

import "testing"

const testAddress = "iot.eclipse.org:1883"

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

func Test_commandCon_receive_ReceiveErr_disconnecting(t *testing.T) {
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

func Test_commandCon_receive_ReceiveErr(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.ctx.wg.Add(1)
	go cmd.receive()

	cmd.ctx.wg.Wait()
}
