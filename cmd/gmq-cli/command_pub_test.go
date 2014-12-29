package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_commandPub_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	cmd, err := newCommandPub([]string{"-t", "topicName"}, cli)
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cmd.run(); err != client.ErrNotYetConnected {
		invalidError(t, err, client.ErrNotYetConnected)
	}
}

func Test_newCommandPub_errCmdArgsParse(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandPub([]string{"-not-exist-flag"}, cli); err != errCmdArgsParse {
		invalidError(t, err, errCmdArgsParse)
	}
}

func Test_newCommandPub(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandPub([]string{"-t", "topicName"}, cli); err != nil {
		nilErrorExpected(t, err)
	}
}
