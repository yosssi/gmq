package main

import "errors"

// Command names
const (
	cmdNameConn = "conn"
	cmdNameHelp = "help"
)

// Multiple error string format
const strErrMulti = "error (%q) occurred while handling the other error (%q)"

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
func newCommand(cmdName string, cmdArgs []string, ctx *context) (command, error) {
	switch cmdName {
	case cmdNameHelp:
		return newCommandHelp(), nil
	case cmdNameConn:
		return newCommandConn(cmdArgs, ctx)
	}

	return nil, errInvalidCmdName
}
