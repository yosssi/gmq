package main

import (
	"errors"

	"github.com/yosssi/gmq/mqtt/client"
)

// Error values
var (
	errInvalidCmdName = errors.New("invalid command name")
	errCmdArgsParse   = errors.New("command arguments parse error")
)

// command represents a command of GMQ Client.
type command interface {
	run() error
}

// newCommand creates and returns a command.
func newCommand(cmdName string, cmdArgs []string, cli *client.Client) (command, error) {
	switch cmdName {
	case cmdNameConn:
		return newCommandConn(cmdArgs, cli)
	case cmdNameDisconn:
		return newCommandDisconn(cli), nil
	case cmdNameHelp:
		return newCommandHelp(), nil
	case cmdNamePub:
		return newCommandPub(cmdArgs, cli)
	case cmdNameQuit:
		return newCommandQuit(cli), nil
	case cmdNameSub:
		return newCommandSub(cmdArgs, cli)
	case cmdNameUnsub:
		return newCommandUnsub(cmdArgs, cli)
	}

	return nil, errInvalidCmdName
}
