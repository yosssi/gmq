package main

import "github.com/yosssi/gmq/mqtt/client"

// Command name
const cmdNameQuit = "quit"

// commandQuit represents a quit command.
type commandQuit struct {
	cli *client.Client
}

// run quits this process.
func (cmd *commandQuit) run() error {
	quit(cmd.cli)

	return nil
}

// newCommandQuit creates and returns a quit command.
func newCommandQuit(cli *client.Client) command {
	return &commandQuit{
		cli: cli,
	}
}

// quit quits this process.
func quit(cli *client.Client) {
	// Disconnect the Network Connection.
	cli.Disconnect()

	// Terminate the Client.
	cli.Terminate()

	// Exit.
	exit(0)
}
