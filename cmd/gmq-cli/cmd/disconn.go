package cmd

import "github.com/yosssi/gmq/mqtt/client"

// Disconn is a command which disconnects the Network Connection to the Server.
var Disconn = &Cmd{
	Name:  "disconn",
	Usage: "disconnect the Network Connection to the Server",
	Run:   disconn,
}

// disconn disconnects the Network Connection to the Server.
func disconn(cli *client.Client, c *Cmd) error {
	return cli.Disconnect()
}
