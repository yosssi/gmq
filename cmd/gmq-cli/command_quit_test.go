package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

const testAddress = "iot.eclipse.org:1883"

func init() {
	exit = func(_ int) {}
}

func Test_commandQuit_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	cmd := commandQuit{
		cli: cli,
	}

	if err := cmd.run(); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_newCommandQuit(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if cmd := newCommandQuit(cli); cmd == nil {
		t.Error("cmd => nil, want => not nil")
	}
}

func Test_quit(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	quit(cli)
}
