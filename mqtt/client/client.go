package client

import (
	"errors"

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
func (cli *Client) Connect(opts *ConnectOptions) error {
	// Return an error if the Client has already connected to the Server.
	if cli.Conn != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnectOptions{}
	}
	opts.Init()

	// Connect to the Server and create a Network Connection.
	conn, err := mqtt.NewConnection(opts.Network, opts.Address)
	if err != nil {
		return err
	}

	// Set the Network Connection to the Client.
	cli.Conn = conn

	return nil
}

// Disconnect closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.Conn == nil {
		return ErrNotYetConnected
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

// SendCONNECT sends a CONNECT Packet to the Server.
func (cli *Client) SendCONNECT(opts *packet.CONNECTOptions) error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.Conn == nil {
		return ErrNotYetConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &packet.CONNECTOptions{}
	}
	opts.Init()

	// Create a Session or reuse the current Session.
	if *opts.CleanSession || cli.Sess == nil {
		// Craete a Session and set it to the Client.
		cli.Sess = NewSession(&SessionOptions{
			CleanSession: opts.CleanSession,
			ClientID:     opts.ClientID,
		})
	} else {
		// Reuse the Session and set its Client Identifier to the Packet options.
		opts.ClientID = cli.Sess.ClientID
	}

	// Create a CONNECT Packet.
	p, err := packet.NewCONNECT(opts)
	if err != nil {
		return err
	}

	// Send the CONNECT Packet to the Server.
	return cli.send(p)
}

// SendDISCONNECT sends a DISCONNECT Packet to the Server.
func (cli *Client) SendDISCONNECT() error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.Conn == nil {
		return ErrNotYetConnected
	}

	// Send a DISCONNECT Packet to the Server.
	return cli.send(packet.NewDISCONNECT())
}

// send sends an MQTT Control Packet to the Server.
func (cli *Client) send(p packet.Packet) error {
	if _, err := p.WriteTo(cli.Conn.W); err != nil {
		return err
	}

	return cli.Conn.W.Flush()
}
