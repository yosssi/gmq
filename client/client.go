package client

import (
	"log"
	"net"
	"strconv"
)

// Client represents an MQTT client.
type Client struct {
	// conns is network connections to the server.
	conns []net.Conn
}

// Conn tries to establish a network connection to the server and
// sends a CONNECT Package to the server.
func (cli *Client) Conn(opts *ConnOpts) error {
	if opts == nil {
		opts = &ConnOpts{}
	}

	opts.Init()

	address := opts.Host + ":" + strconv.Itoa(int(*opts.Port))

	log.Printf("Connecting to %s", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("Connected successfully to %s", address)

	cli.conns = append(cli.conns, conn)

	return nil
}

// New creates and returns an MQTT client.
func New() *Client {
	return &Client{}
}
