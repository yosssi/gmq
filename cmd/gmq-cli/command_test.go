package main

import (
	"testing"

	"github.com/yosssi/gmq/mqtt/client"
)

func Test_newCommand(t *testing.T) {
	cmdNames := []string{
		cmdNameConn,
		cmdNameDisconn,
		cmdNameHelp,
		cmdNamePub,
		cmdNameQuit,
		cmdNameSub,
		cmdNameUnsub,
	}

	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	for _, cmdName := range cmdNames {
		if _, err := newCommand(cmdName, nil, cli); err != nil {
			nilErrorExpected(t, err)
		}
	}
}

func Test_newCommand_errInvalidCmdName(t *testing.T) {
	cli := client.New(&client.Options{
		ErrorHandler: func(_ error) {},
	})

	defer quit(cli)

	if _, err := newCommand("invalidCmdName", nil, cli); err != errInvalidCmdName {
		invalidError(t, err, errInvalidCmdName)
	}
}
