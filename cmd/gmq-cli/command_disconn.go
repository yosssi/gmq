package main

import "github.com/yosssi/gmq/mqtt/client"

// Command name
const cmdNameDisconn = "disconn"

// commandDisconn represents a disconn command.
type commandDisconn struct {
	cli *client.Client
}

// run sends a DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cmd *commandDisconn) run() error {
	return cmd.cli.Disconnect()
}

// newCommandDisconn creates and returns a disconn command.
func newCommandDisconn(cli *client.Client) command {
	return &commandDisconn{
		cli: cli,
	}
}
