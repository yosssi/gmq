package common

import "net"

// Connection represents a Network Connection.
type Connection struct {
	net.Conn
}

// NewConnection connects to the address on the named network,
// creates a Connection and returns it.
func NewConnection(network, address string) (*Connection, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return &Connection{Conn: conn}, nil
}
