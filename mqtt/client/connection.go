package client

import (
	"bufio"
	"net"
)

// connection represents a Network Connection.
type connection struct {
	net.Conn
	// r is the buffered reader.
	r *bufio.Reader
	// w is the buffered writer.
	w *bufio.Writer
	// disconnected is true if the Network Connection
	// has been disconnected.
	disconnected bool
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
