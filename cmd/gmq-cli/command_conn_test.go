package main

import (
	"path/filepath"
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_commandConn_run(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	cmd, err := newCommandConn(nil, cli)
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if err := cmd.run(); err == nil {
		notNilErrorExpected(t)
	}
}

func Test_newCommandConn_errCmdArgsParse(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandConn([]string{"-not-exist-flag"}, cli); err != errCmdArgsParse {
		invalidError(t, err, errCmdArgsParse)
	}
}

func Test_newCommandConn_ReadFileErr(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandConn([]string{"-crt", "not_exist_file.crt"}, cli); err == nil {
		notNilErrorExpected(t)
	}
}

func Test_newCommandConn_errParseCrtFailure(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandConn([]string{"-crt", filepath.Join("test", "error.crt")}, cli); err != errParseCrtFailure {
		invalidError(t, err, errParseCrtFailure)
	}
}

func Test_newCommandConn(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommandConn([]string{"-crt", filepath.Join("test", "test.crt")}, cli); err != nil {
		nilErrorExpected(t, err)
	}
}

func notNilErrorExpected(t *testing.T) {
	t.Error("err => nil, want => not nil")
}
