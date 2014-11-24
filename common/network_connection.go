package common

import "net"

// NetworkConnection represents a Network Connection.
type NetworkConnection struct {
	conn net.Conn
}

// NewNetworkConnection connects to the address on the named network,
// creates a Network Connection and returns it.
func NewNetworkConnection(network, address string) (*NetworkConnection, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &NetworkConnection{conn: conn}, nil
}
