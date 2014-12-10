package main

import (
	"errors"

	"github.com/yosssi/gmq/mqtt/client"
)

// Command names
const (
	cmdNameConn = "conn"
	cmdNameHelp = "help"
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
	case cmdNameHelp:
		return newCommandHelp(), nil
	case cmdNameConn:
		return newCommandConn(cmdArgs, cli)
	}

	return nil, errInvalidCmdName
}
