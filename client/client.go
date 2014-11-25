package client

import (
	"errors"
	"sync"

	"github.com/yosssi/gmq/common"
)

// Error values
var ErrAlreadyConnected = errors.New("the Client has already connected to the Server")

// Client represents a Client.
type Client struct {
	// mu is a reader/writer mutual exclusion lock for the Client.
	mu sync.RWMutex
	// networkConnection is a Network Connection.
	conn *common.Connection
}

// Connect tries to establish a network connection to the Server and
// sends a CONNECT Package to the Server.
func (cli *Client) Connect(address string, opts *common.OptionsPacketCONNECT) error {
	// Lock for the update of the Client's fields.
	cli.mu.Lock()
	defer cli.mu.Unlock()

	// Return an error if the Client has already connected to the Server.
	if cli.conn != nil {
		return ErrAlreadyConnected
	}

	// Connect to the Server and create a Network Connection.
	conn, err := common.NewConnection("tcp", address)
	if err != nil {
		return err
	}
	cli.conn = conn

	return nil
}

// New creates and returns a Client.
func New() *Client {
	return &Client{}
}
