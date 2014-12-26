package client

import (
	"bufio"
	"net"
	"sync"

	"github.com/yosssi/gmq/mqtt/packet"
)

// Buffer size of the send channel
const sendBufSize = 1024

// connection represents a Network Connection.
type connection struct {
	net.Conn
	// r is the buffered reader.
	r *bufio.Reader
	// w is the buffered writer.
	w *bufio.Writer
	// disconnected is true if the Network Connection
	// has been disconnected by the Client.
	disconnected bool

	// wg is the Wait Group for the goroutines
	// which are launched by the Connect method.
	wg sync.WaitGroup
	// connack is the channel which handles the signal
	// to notify the arrival of the CONNACK Packet.
	connack chan struct{}
	// send is the channel which handles the Packet.
	send chan packet.Packet
	// sendEnd is the channel which ends the goroutine
	// which sends a Packet to the Server.
	sendEnd chan struct{}

	// muPINGRESPs is the Mutex for pingresps.
	muPINGRESPs sync.RWMutex
	// pingresps is the slice of the channels which
	// handle the signal to notify the arrival of
	// the PINGRESP Packet.
	pingresps []chan struct{}

	// subscriptions contains the subscription information.
	subscriptions map[string]*SubState
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
		Conn:          conn,
		r:             bufio.NewReader(conn),
		w:             bufio.NewWriter(conn),
		connack:       make(chan struct{}, 1),
		send:          make(chan packet.Packet, sendBufSize),
		sendEnd:       make(chan struct{}, 1),
		subscriptions: make(map[string]*SubState),
	}

	// Return the Network Connection.
	return c, nil
}
