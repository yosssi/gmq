package client

import "time"

// ConnectOptions represents options for the Connect method
// of the Client.
type ConnectOptions struct {
	// Network is the network on which the Client connects to.
	Network string
	// Address is the address which the Client connects to.
	Address string
	// CONNACKTimeout is timeout in seconds for the Client
	// to wait for receiving the CONNACK Packet after sending
	// the CONNECT Packet.
	CONNACKTimeout time.Duration
}
