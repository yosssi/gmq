package client

import (
	"errors"
	"log"
	"net"
	"strconv"
	"sync"
)

// Error values
var ErrAlreadyConnected = errors.New("the client has already connected to the server")

// Client represents an MQTT client.
type Client struct {
	sync.RWMutex
	// conns is network connections to the server.
	conn net.Conn
}

// Conn tries to establish a network connection to the server and
// sends a CONNECT Package to the server.
func (cli *Client) Conn(opts *ConnOpts) error {
	cli.Lock()
	defer cli.Unlock()

	if cli.conn != nil {
		return ErrAlreadyConnected
	}

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

	cli.conn = conn

	return nil
}

// New creates and returns an MQTT client.
func New() *Client {
	return &Client{}
}
