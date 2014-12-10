package main

import (
	"errors"

	"github.com/yosssi/gmq/mqtt/client"
)

// Command names
const (
	cmdNameHelp = "help"
)

// Error values
var (
	errInvalidCmdName = errors.New("invalid command name")
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
	}

	return nil, errInvalidCmdName
}
