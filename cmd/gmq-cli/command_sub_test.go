package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_commandSub_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	cmd, err := newCommandSub([]string{"-t", "topicFilter"}, cli)
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cmd.run(); err != client.ErrNotYetConnected {
		invalidError(t, err, client.ErrNotYetConnected)
	}
}

func Test_newCommandSub_errCmdArgsParse(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandSub([]string{"-not-exist-flag"}, cli); err != errCmdArgsParse {
		invalidError(t, err, errCmdArgsParse)
	}
}

func Test_newCommandSub(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandSub([]string{"-t", "topicFilter"}, cli); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_messageHandler(t *testing.T) {
	messageHandler([]byte("topicName"), []byte("message"))
}
