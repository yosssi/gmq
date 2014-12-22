package client

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/packet"
)

// Client represents a Client.
type Client struct {
	// mu is the Mutex for the Network Connection.
	mu sync.RWMutex
	// conn is the Network Connection.
	conn *connection
	// sess is the Session.
	sess *session
}

// Connect establishes a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(opts *ConnectOptions, packetOpts *packet.CONNECTOptions) error {
	return nil
}
