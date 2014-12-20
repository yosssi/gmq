package client

import (
	"errors"
	"fmt"
	"io"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Multiple error string format
const strErrMulti = "error (%q) occurred while handling the other error (%q)"

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
)

// closeConn calls the Client's close method.
// This global variable is defined to make writing tests easy and
// another function value will be assigned to this global variable
// while testing.
var closeConn = func(cli *Client) error {
	return cli.Close()
}

// sendDISCONNECT calls the Client's sendDISCONNECT method.
// This global variable is defined to make writing tests easy and
// another function value will be assigned to this global variable
// while testing.
var sendDISCONNECT = func(cli *Client) error {
	return cli.Send(packet.NewDISCONNECT())
}

// Client represents a Client.
type Client struct {
	// conn is the Network Connection.
	conn *mqtt.Connection
	// Sess is the Session.
	sess *Session
}

// Connect tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(network, address string, packetOpts *packet.CONNECTOptions) error {
	// Try to establish a Network Connection to the Server.
	if err := cli.establish(network, address); err != nil {
		return err
	}

	// Send a CONNECT Packet to the Server.
	if err := cli.sendCONNECT(packetOpts); err != nil {
		// Close the Network Connection to the Server.
		if anotherErr := closeConn(cli); anotherErr != nil {
			return fmt.Errorf(strErrMulti, anotherErr, err)
		}

		// Clear the Session.
		cli.ClearConnection()

		return err
	}

	return nil
}

// Disconnect sends a DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Send a DISCONNECT Packet to the Server.
	if err := sendDISCONNECT(cli); err != nil {
		return err
	}

	// Close the Network Connection.
	if err := closeConn(cli); err != nil {
		return err
	}

	return nil
}

// ClearConnection clears the Network Connection.
func (cli *Client) ClearConnection() {
	// Clear the Network Connection of the Client.
	cli.conn = nil

	// Clear the Session if the CleanSession is true.
	if cli.sess != nil && cli.sess.CleanSession {
		cli.sess = nil
	}
}

// Send sends an MQTT Control Packet to the Server.
func (cli *Client) Send(p packet.Packet) error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Write the Packet to the buffered writer.
	if _, err := p.WriteTo(cli.conn.W); err != nil {
		return err
	}

	// Flush the buffered writer.
	return cli.conn.W.Flush()
}

// Receive receives an MQTT Control Packet from the Server.
func (cli *Client) Receive() (packet.Packet, error) {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return nil, ErrNotYetConnected
	}

	// Get the first byte of the Packet.
	b, err := cli.conn.R.ReadByte()
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
		b, err = cli.conn.R.ReadByte()
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
		if _, err = io.ReadFull(cli.conn.R, remaining); err != nil {
			return nil, err
		}
	}

	// Create and return a Packet.
	return packet.NewFromBytes(fixedHeader, remaining)
}

// Close closes the Network Connection.
func (cli *Client) Close() error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Close the Network Connection.
	if err := cli.conn.Close(); err != nil {
		return err
	}

	return nil
}

// Connected returns true if the Client has a Network Connection.
func (cli *Client) Connected() bool {
	return cli.conn != nil
}

// establish tries to establish a Network Connection to the Server.
func (cli *Client) establish(network, address string) error {
	// Return an error if the Client has already connected to the Server.
	if cli.conn != nil {
		return ErrAlreadyConnected
	}

	// Connect to the Server and create a Network Connection.
	conn, err := mqtt.NewConnection(network, address)
	if err != nil {
		return err
	}

	// Set the Network Connection to the Client.
	cli.conn = conn

	return nil
}

// SendCONNECT sends a CONNECT Packet to the Server.
func (cli *Client) sendCONNECT(opts *packet.CONNECTOptions) error {
	// Initialize the options.
	if opts == nil {
		opts = &packet.CONNECTOptions{}
	}
	opts.Init()

	// Create a Session or reuse the current Session.
	if *opts.CleanSession || cli.sess == nil {
		// Craete a Session and set it to the Client.
		cli.sess = NewSession(&SessionOptions{
			CleanSession: opts.CleanSession,
			ClientID:     opts.ClientID,
		})
	} else {
		// Reuse the Session and set its Client Identifier to the Packet options.
		opts.ClientID = cli.sess.ClientID
	}

	// Create a CONNECT Packet.
	p, err := packet.NewCONNECT(opts)
	if err != nil {
		return err
	}

	// Send the CONNECT Packet to the Server.
	return cli.Send(p)
}

// New creates and returns a Client.
func New(_ *Options) *Client {
	return &Client{}
}
