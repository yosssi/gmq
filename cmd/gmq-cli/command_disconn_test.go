package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_commandDisconn_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	cmd := &commandDisconn{cli: cli}

	if err := cmd.run(); err != client.ErrNotYetConnected {
		invalidError(t, err, client.ErrNotYetConnected)
	}
}

func Test_newCommandDisconn(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if cmd := newCommandDisconn(cli); cmd == nil {
		t.Error("cmd => nil, want => not nil")
	}
}
