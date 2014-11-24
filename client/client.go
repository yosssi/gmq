package client

import (
	"errors"
	"strconv"
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
	networkConnection *common.NetworkConnection
}

// Conn tries to establish a network connection to the Server and
// sends a CONNECT Package to the Server.
func (cli *Client) Conn(opts *ConnOpts) error {
	// Lock for the update of the Client's field.
	cli.mu.Lock()
	defer cli.mu.Unlock()

	// Return an error if the Client has already connected to the Server.
	if cli.networkConnection != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnOpts{}
	}
	opts.Init()

	// Connect to the Server and create a Network Connection.
	networkConnection, err := common.NewNetworkConnection("tcp", opts.Host+":"+strconv.Itoa(int(*opts.Port)))
	if err != nil {
		return err
	}
	cli.networkConnection = networkConnection

	return nil
}

// New creates and returns a Client.
func New() *Client {
	return &Client{}
}
