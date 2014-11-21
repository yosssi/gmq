package client

import "net"

// Client represents an MQTT client.
type Client struct {
	// conns is network connections to the server.
	conns []net.Conn
}

// New creates and returns an MQTT client.
func New() *Client {
	return &Client{}
}
