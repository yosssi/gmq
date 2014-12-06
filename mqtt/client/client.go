package client

import (
	"bufio"
	"errors"
	"io"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
)

// Client represents a Client.
type Client struct {
	// Conn is the Network Connection.
	Conn *mqtt.Connection
	// Sess is the Session.
	Sess *Session
}

// Connect tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(opts *ConnectOptions, packetOpts *packet.CONNECTOptions) error {
	// Return an error if the Client has already connected to the Server.
	if cli.Conn != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnectOptions{}
	}
	opts.Init()

	// Initialize the options for the CONNECT Packet.
	if packetOpts == nil {
		packetOpts = &packet.CONNECTOptions{}
	}
	packetOpts.Init()

	// Connect to the Server and create a Network Connection.
	conn, err := mqtt.NewConnection(opts.Network, opts.Address)
	if err != nil {
		return err
	}
	cli.Conn = conn

	// Create a Session or reuse the current Session.
	if *packetOpts.CleanSession || cli.Sess == nil || cli.Sess.CleanSession {
		// Craete a Session and set it to the Client.
		cli.Sess = NewSession(&SessionOptions{
			CleanSession: packetOpts.CleanSession,
			ClientID:     packetOpts.ClientID,
		})
	} else {
		// Reuse the Session and set its Client Identifier to the Packet options.
		packetOpts.ClientID = cli.Sess.ClientID
	}

	// Create a CONNECT Packet.
	p, err := packet.NewCONNECT(packetOpts)
	if err != nil {
		return err
	}

	// Send the CONNECT Packet to the Server.
	return cli.Send(p)
}

// Disconnect sends the DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.Conn == nil {
		return ErrNotYetConnected
	}

	// Send the DISCONNECT Packet to the Server.
	if err := cli.Send(packet.NewDISCONNECT()); err != nil {
		return err
	}

	// Close the Network Connection.
	if err := cli.Conn.Close(); err != nil {
		return err
	}

	// Clear the Network Connection of the Client.
	cli.Conn = nil

	// Clear the Session if the CleanSession is true.
	if cli.Sess != nil && cli.Sess.CleanSession {
		cli.Sess = nil
	}

	return nil
}

// Send sends an MQTT Control Packet to the Server.
func (cli *Client) Send(p packet.Packet) error {
	if _, err := p.WriteTo(cli.Conn.W); err != nil {
		return err
	}

	return cli.Conn.W.Flush()
}

// Receive receives an MQTT Control Packet from the Server.
func Receive(r *bufio.Reader) (byte, packet.Packet, error) {
	// Get the first byte of the Packet.
	b, err := r.ReadByte()
	if err != nil {
		return 0x00, nil, err
	}

	// Extract the MQTT Control Packet Type from the first byte.
	ptype := b >> 4

	// Create the Fixed header.
	fixedHeader := []byte{b}

	// Get and decode the Remaining Length.
	var mp uint32 = 1 // multiplier
	var rl uint32     // the Remaining Length
	for {
		b, err = r.ReadByte()
		if err != nil {
			return 0x00, nil, err
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
		if _, err = io.ReadFull(r, remaining); err != nil {
			return 0x00, nil, err
		}
	}

	var p packet.Packet

	switch ptype {
	case packet.TypeCONNACK:
		// Create the CONNACK Packet from the byte data to validate the data.
		if p, err = packet.NewCONNACKFromBytes(fixedHeader, remaining); err != nil {
			return 0x00, nil, err
		}
	}

	return ptype, p, nil
}
