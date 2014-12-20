package main

import "errors"

// Command names
const (
	cmdNameConn    = "conn"
	cmdNameDisconn = "disconn"
	cmdNameHelp    = "help"
	cmdNamePub     = "pub"
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
func newCommand(cmdName string, cmdArgs []string, ctx *context) (command, error) {
	switch cmdName {
	case cmdNameConn:
		return newCommandConn(cmdArgs, ctx)
	case cmdNameDisconn:
		return newCommandDisconn(ctx), nil
	case cmdNameHelp:
		return newCommandHelp(), nil
	case cmdNamePub:
		return newCommandPub(cmdArgs, ctx)
	}

	return nil, errInvalidCmdName
}
