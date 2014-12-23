package client

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/yosssi/gmq/mqtt/packet"
)

// Multiple errors string format
const strErrMulti = "error (%q) occurred while handling the other error (%q)"

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
)

// Client represents a Client.
type Client struct {
	// mu is the Mutex for the Network Connection.
	mu sync.RWMutex
	// conn is the Network Connection.
	conn *connection
	// sess is the Session.
	sess *session

	// wg is the Wait Group for the goroutines.
	wg sync.WaitGroup
}

// Connect establishes a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(opts *ConnectOptions) error {
	// Lock for the connection.
	cli.mu.Lock()

	// Unlock.
	defer cli.mu.Unlock()

	// Return an error if the Client has already connected to the Server.
	if cli.conn != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnectOptions{}
	}

	// Establish a Network Connection.
	conn, err := newConnection(opts.Network, opts.Address)
	if err != nil {
		return err
	}

	// Set the Network Connection to the Client.
	cli.conn = conn

	// Create a Session or reuse the current Session.
	if opts.CleanSession || cli.sess == nil {
		// Create a Session and set it to the Client.
		cli.sess = &session{
			cleanSession: opts.CleanSession,
			clientID:     opts.ClientID,
		}
	} else {
		// Reuse the Session and set its Client Identifier to the options.
		opts.ClientID = cli.sess.clientID
	}

	// Send a CONNECT Packet to the Server.
	err = cli.sendCONNECT(&packet.CONNECTOptions{
		ClientID:     opts.ClientID,
		UserName:     opts.UserName,
		Password:     opts.Password,
		CleanSession: opts.CleanSession,
		KeepAlive:    opts.KeepAlive,
		WillTopic:    opts.WillTopic,
		WillMessage:  opts.WillMessage,
		WillQoS:      opts.WillQoS,
		WillRetain:   opts.WillRetain,
	})

	if err != nil {
		// Close the Network Connection.
		if anotherErr := cli.conn.Close(); anotherErr != nil {
			return fmt.Errorf(strErrMulti, anotherErr)
		}

		// Clean the Network Connection and the Session if necessary.
		cli.clean()

		return err
	}

	// Launch a goroutine which receives a Packet from the Server.
	cli.wg.Add(1)
	func() {
		defer cli.wg.Done()

		for {
			cli.receive()
		}
	}()

	return nil
}

// send sends an MQTT Control Packet to the Server.
func (cli *Client) send(p packet.Packet) error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Write the Packet to the buffered writer.
	if _, err := p.WriteTo(cli.conn.w); err != nil {
		return err
	}

	// Flush the buffered writer.
	return cli.conn.w.Flush()
}

// sendCONNECT creates a CONNECT Packet and sends it to the Server.
func (cli *Client) sendCONNECT(opts *packet.CONNECTOptions) error {
	// Initialize the options.
	if opts == nil {
		opts = &packet.CONNECTOptions{}
	}

	// Create a CONNECT Packet.
	p, err := packet.NewCONNECT(opts)
	if err != nil {
		return err
	}

	// Send a CONNECT Packet to the Server.
	return cli.send(p)
}

// receive receives an MQTT Control Packet from the Server.
func (cli *Client) receive() (packet.Packet, error) {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return nil, ErrNotYetConnected
	}

	// Get the first byte of the Packet.
	b, err := cli.conn.r.ReadByte()
	if err != nil {
		return nil, err
	}

	// Create the Fixed header.
	fixedHeader := []byte{b}

	// Get and decode the Remaining Length.
	var mp uint32 = 1 // multiplier
	var rl uint32     // the Remaining Length
	for {
		// Get the next byte of the Packet.
		b, err = cli.conn.r.ReadByte()
		if err != nil {
			return nil, err
		}

		fixedHeader = append(fixedHeader, b)

		rl += uint32(b&0x7F) * mp

		if b&0x80 == 0 {
			break
		}

		mp *= 128
	}

	// Create the Remaining (the Variable header and the Payload).
	remaining := make([]byte, rl)

	if rl > 0 {
		// Get the remaining of the Packet.
		if _, err = io.ReadFull(cli.conn.r, remaining); err != nil {
			return nil, err
		}
	}

	// Create and return a Packet.
	fmt.Println(fixedHeader, remaining)
	return nil, nil
	//return packet.NewFromBytes(fixedHeader, remaining)
}

// clean cleans the Network Connection and the Session if necessary.
func (cli *Client) clean() {
	// Clean the Network Connection.
	cli.conn = nil

	// Clean the Session if the Clean Session is true.
	if cli.sess != nil && cli.sess.cleanSession {
		cli.sess = nil
	}
}
