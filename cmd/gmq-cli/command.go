package main

import "github.com/yosssi/gmq/mqtt/client"

type command interface {
	run() error
}

// newCommand creates and returns a command.
func newCommand(cmdName string, cmdArgs []string, cli *client.Client) command {
	return nil
}
