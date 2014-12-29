package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_newCommandUnsub_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	cmd, err := newCommandUnsub([]string{"-t", "topicName"}, cli)
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cmd.run(); err != client.ErrNotYetConnected {
		invalidError(t, err, client.ErrNotYetConnected)
	}
}

func Test_newCommandUnsub_errCmdArgsParse(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandUnsub([]string{"-not-exist-flag"}, cli); err != errCmdArgsParse {
		invalidError(t, err, errCmdArgsParse)
	}
}

func Test_newCommandUnsub(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandUnsub([]string{"-t", "topicName"}, cli); err != nil {
		nilErrorExpected(t, err)
	}
}

func invalidError(t *testing.T, err, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
