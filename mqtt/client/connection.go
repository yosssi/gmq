package client

import (
	"bufio"
	"net"
)

// connection represents a Network Connection.
type connection struct {
	net.Conn
	r *bufio.Reader
	w *bufio.Writer
}

// newConnection connects to the address on the named network,
// creates a Network Connection and returns it.
func newConnection(network, address string) (*connection, error) {
	// Connect to the address on the named network.
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	// Create a Network Connection.
	c := &connection{
		Conn: conn,
		r:    bufio.NewReader(conn),
		w:    bufio.NewWriter(conn),
	}

	// Return the Network Connection.
	return c, nil
}
