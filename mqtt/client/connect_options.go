package client

import "time"

// Default values
const (
	DefaultNetwork = "tcp"
	DefaultAddress = "localhost:1883"
)

// Default values
var (
	DefaultCONNACKTimeout = 60 * time.Second
)

// ConnectOptions is options for the Client's Connect method.
type ConnectOptions struct {
	// Network is the network which the Client connects to.
	Network string
	// Address is the address which the Client connects to.
	Address string
	// CONNACKTimeout is the timeout for waiting for receiving the CONNACK Packet
	// sent from the Server.
	CONNACKTimeout time.Duration
}

// Init initializes the ConnectOptions.
func (opts *ConnectOptions) Init() {
	if opts.Network == "" {
		opts.Network = DefaultNetwork
	}

	if opts.Address == "" {
		opts.Address = DefaultAddress
	}

	if opts.CONNACKTimeout == 0 {
		opts.CONNACKTimeout = DefaultCONNACKTimeout
	}
}
