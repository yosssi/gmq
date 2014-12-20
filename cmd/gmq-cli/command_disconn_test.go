package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_commandDisconn_run(t *testing.T) {
	cmd := newCommandDisconn(newContext())
	if err := cmd.run(); err != client.ErrNotYetConnected {
		errorfErr(t, err, client.ErrNotYetConnected)
	}
}

func Test_disconnect(t *testing.T) {
	ctx := newContext()

	if err := ctx.cli.Connect("tcp", testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := ctx.cli.Close(); err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if err := disconnect(ctx); err == nil {
		t.Error("err => nil, want => not nil")
	}
}
