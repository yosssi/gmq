package client

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/yosssi/gmq/common"
	"github.com/yosssi/gmq/common/packet"
)

// Defalut values
const (
	defaultErrcBufferSize  = 1024
	defaultSendcBufferSize = 1024
)

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
	ErrNotCONNACK       = errors.New("the Packet which was not the CONNACK Packet has been received")
	ErrCONNACKTimeout   = errors.New("Timeout has occurred while waiting for receiving the CONNACK Packet from the Server")
)

// Client represents a Client.
type Client struct {
	// Errc is a channel handling errors which are sent by the goroutines
	// which sends or receives MQTT Control Packets.
	Errc chan error
	// mu is a reader/writer mutual exclusion lock for the Client.
	mu sync.RWMutex
	// networkConnection is a Network Connection.
	conn *common.Connection
	// sendc is a channel handling MQTT Control Packets which are sent from
	// the Client to the Server.
	sendc chan packet.Packet
}

// Connect tries to establish a network connection to the Server and
// sends a CONNECT Package to the Server.
func (cli *Client) Connect(opts *ConnectOptions, packetOpts *packet.CONNECTOptions) error {
	// Lock for the update of the Client's fields.
	cli.mu.Lock()
	defer cli.mu.Unlock()

	// Return an error if the Client has already connected to the Server.
	if cli.conn != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnectOptions{}
	}
	opts.Init()

	if packetOpts == nil {
		packetOpts = &packet.CONNECTOptions{}
	}
	packetOpts.Init()

	// Connect to the Server and create a Network Connection.
	conn, err := common.NewConnection(opts.Network, opts.Address)
	if err != nil {
		return err
	}
	cli.conn = conn

	// Send the CONNECT Packet to the Server.
	if err := cli.send(packet.NewCONNECT(packetOpts)); err != nil {
		// TODO disconnect
		return err
	}

	// Wait for receiving the CONNACK Packet.
	connacked := make(chan struct{})

	go func() {
		p, err := cli.receive()
		if err != nil {
			cli.Errc <- err
		}

		if _, ok := p.(*packet.CONNACK); !ok {
			cli.Errc <- ErrNotCONNACK
		}

		connacked <- struct{}{}
	}()

	select {
	case <-connacked:
	case err := <-cli.Errc:
		// TODO disconnect
		return err
	case <-time.After(opts.CONNACKTimeout):
		// TODO disconnect
		return ErrCONNACKTimeout
	}

	// Create a send channel handling MQTT Control Packets and set it to the Client.
	cli.sendc = make(chan packet.Packet, defaultSendcBufferSize)

	// Launch a goroutine which sends MQTT Control Packets to the Server.
	go func() {
		// Send MQTT Control Packets.
		for p := range cli.sendc {
			if err := cli.send(p); err != nil {
				// TODO disconnect
				cli.Errc <- err
				break
			}
		}
	}()

	// Launch a goroutine which receives MQTT Control Packets from the Server.
	go func() {
		// Receive MQTT Control Packets from the Server.
		for {
			if _, err := cli.receive(); err != nil {
				// TODO disconnect
				cli.Errc <- err
				break
			}
		}
	}()

	return nil
}

// Disconnect sends the DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Lock for the update of the Client's fields.
	cli.mu.Lock()
	defer cli.mu.Unlock()

	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Send the DISCONNECT Packet to the Server.

	// Close the Network Connection.
	return cli.conn.Close()
}

// send sends an MQTT Control Packet to the Server.
func (cli *Client) send(p packet.Packet) error {
	if _, err := p.WriteTo(cli.conn.W); err != nil {
		return err
	}

	return cli.conn.W.Flush()
}

// receive receives MQTT Control Packets from the Server
func (cli *Client) receive() (packet.Packet, error) {
	// Get the first byte of the Packet.
	b, err := cli.conn.R.ReadByte()
	if err != nil {
		return nil, err
	}

	// Extract the MQTT Control Packet Type from the first byte.
	packetType := b >> 4

	// Create the Fixed header.
	fixedHeader := []byte{b}

	// Get and decode the Remaining Length.
	var mp uint32 = 1 // multiplier
	var rl uint32     // the Remaining Length
	for {
		b, err = cli.conn.R.ReadByte()
		if err != nil {
			return nil, err
		}

		fixedHeader = append(fixedHeader, b)

		rl += uint32(b&127) * mp

		if b&128 == 0 {
			break
		}

		mp *= 128
	}

	// Create the Remaining (the Variable header and the Payload).
	remaining := make([]byte, rl)

	if rl > 0 {
		if _, err = io.ReadFull(cli.conn.R, remaining); err != nil {
			return nil, err
		}
	}

	var p packet.Packet

	switch packetType {
	case packet.TypeCONNACK:
		// Create the CONNACK Packet from the byte data to validate the data.
		if p, err = packet.NewCONNACKFromBytes(fixedHeader, remaining); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// New creates and returns a Client.
func New() *Client {
	return &Client{
		Errc: make(chan error, defaultErrcBufferSize),
	}
}
