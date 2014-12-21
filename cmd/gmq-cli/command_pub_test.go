package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

func Test_commandPub_run_errNotYetConnected(t *testing.T) {
	cmd, err := newCommand(cmdNamePub, nil, newContext())
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cmd.run(); err != client.ErrNotYetConnected {
		errorfErr(t, err, client.ErrNotYetConnected)
	}
}

func Test_commandPub_run_newPUBLISHErr(t *testing.T) {
	ctx := newContext()

	if err := ctx.cli.Connect("tcp", testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cmd, err := newCommand(cmdNamePub, nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cmd.(*commandPub).publishOpts = &packet.PUBLISHOptions{
		QoS: 3,
	}

	if err := cmd.run(); err != packet.ErrPUBLISHInvalidQoS {
		errorfErr(t, err, packet.ErrPUBLISHInvalidQoS)
		return
	}

	if err := ctx.cli.Disconnect(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_commandPub_run(t *testing.T) {
	ctx := newContext()

	if err := ctx.cli.Connect("tcp", testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	cmd, err := newCommand(cmdNamePub, nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := cmd.run(); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := ctx.cli.Disconnect(); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}
}
