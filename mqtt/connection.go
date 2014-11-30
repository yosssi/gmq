package mqtt

import (
	"bufio"
	"net"
)

// Connection represents a Network Connection.
type Connection struct {
	net.Conn
	R *bufio.Reader
	W *bufio.Writer
}

// NewConnection connects to the address on the named network,
// creates a Network Connection and returns it.
func NewConnection(network, address string) (*Connection, error) {
	// Connect to the address on the named network
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	// Create a Network Connection.
	c := &Connection{
		Conn: conn,
		R:    bufio.NewReader(conn),
		W:    bufio.NewWriter(conn),
	}

	return c, nil
}
