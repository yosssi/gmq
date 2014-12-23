package client

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt/packet"
)

// Buffer size of the send channel
const sendBufSize = 1024

// Multiple errors string format
const strErrMulti = "error (%q) occurred while handling the other error (%q)"

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
	ErrCONNACKTimeout   = errors.New("the CONNACK Packet was not received within a reasonalbe amount of time")
)

// Client represents a Client.
type Client struct {
	// mu is the Mutex for the Network Connection.
	mu sync.RWMutex
	// conn is the Network Connection.
	conn *connection
	// sess is the Session.
	sess *session

	// wgNew is the Wait Group for the goroutines
	// which are launched by the New method.
	wgNew sync.WaitGroup
	// disconnc is the channel which handles the signal
	// to disconnect the Network Connection.
	disconnc chan struct{}
	// disconnEndc is the channel which ends the goroutine
	// which disconnects the Network Connection.
	disconnEndc chan struct{}

	// wgConn is the Wait Group for the goroutines
	// which are launched by the Connect method.
	wgConn sync.WaitGroup
	// connackc is the channel which handles the signal
	// to notify the arrival of the CONNACK Packet.
	connackc chan struct{}
	// connackEndc is the channel which ends the goroutine
	// which monitors the arrival of the CONNACK Packet.
	connackEndc chan struct{}
	// sendc is the channel which handles the Packet.
	sendc chan packet.Packet
	// sendEndc is the channel which ends the goroutine
	// which sends a Packet to the Server.
	sendEndc chan struct{}

	// errHandler is the error handler.
	errHandler func(error)
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

	// Launch a goroutine which waits for receiving the CONNACK Packet.
	cli.wgConn.Add(1)
	go cli.waitCONNACK(opts.CONNACKTimeout)

	// Launch a goroutine which receives a Packet from the Server.
	cli.wgConn.Add(1)
	go cli.receivePackets()

	// Launch a goroutine which sends a Packet to the Server.
	cli.wgConn.Add(1)
	go cli.sendPackets(opts.KeepAlive)

	return nil
}

// Disconnect sends a DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Lock for the disconnection.
	cli.mu.Lock()

	// Return an error if the Client has not yet connected to the Server.
	if cli.conn != nil {
		// Unlock.
		cli.mu.Unlock()

		return ErrNotYetConnected
	}

	// Send a DISCONNECT Packet to the Server.
	// Ignore the error returned by the send method because
	// we proceed to the subsequent disconnecting processing
	// even if the send method returns the error.
	cli.send(packet.NewDISCONNECT())

	// Close the Network Connection.
	if err := cli.conn.Close(); err != nil {
		// Unlock.
		cli.mu.Unlock()

		return err
	}

	// Change the state of the Network Connection to disconnected.
	cli.conn.disconnected = true

	// Unlock.
	cli.mu.Unlock()

	// Send the end signals to the goroutines via the channels.
	select {
	case cli.connackEndc <- struct{}{}:
	default:
	}

	select {
	case cli.sendEndc <- struct{}{}:
	default:
	}

	// Wait until all goroutines end.
	cli.wgConn.Wait()

	// Lock for the cleaning of the Network Connection and the Session.
	cli.mu.Lock()

	// Clean the Network Connection and the Session.
	cli.clean()

	// Initialize the channels of the Client.
	cli.initChans()

	// Unlock.
	cli.mu.Unlock()

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
	fixedHeader := packet.FixedHeader([]byte{b})

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
	return packet.NewFromBytes(fixedHeader, remaining)
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

// initChans initializes the channels of the client.
func (cli *Client) initChans() {
	cli.connackc = make(chan struct{}, 1)
	cli.connackEndc = make(chan struct{}, 1)
	cli.sendc = make(chan packet.Packet, sendBufSize)
	cli.sendEndc = make(chan struct{}, 1)
}

// waitCONNACK waits for receiving the CONNACK Packet.
func (cli *Client) waitCONNACK(timeout time.Duration) {
	defer cli.wgConn.Done()

	var timeoutc <-chan time.Time

	if timeout > 0 {
		timeoutc = time.After(timeout * time.Second)
	}

	select {
	case <-cli.connackc:
	case <-timeoutc:
		// Handle the timeout error.
		if cli.errHandler != nil {
			cli.errHandler(ErrCONNACKTimeout)
		}

		// Sned a disconnect siganl to the goroutine
		// via the channel if possible.
		select {
		case cli.disconnc <- struct{}{}:
		default:
		}
	case <-cli.connackEndc:
	}
}

// receivePackets receives Packets from the Server.
func (cli *Client) receivePackets() {
	defer cli.wgConn.Done()

	for {
		// Receive a Packet from the Server.
		p, err := cli.receive()
		if err != nil {
			// Handle the error and disconnect
			// the Network Connection.
			cli.handleErrorAndDisconn(err)

			// End the goroutine.
			return
		}

		// Handle the Packet.
		if err := cli.handlePacket(p); err != nil {
			fmt.Println(err)
		}
	}
}

// handlePacket handles the Packet.
func (cli *Client) handlePacket(p packet.Packet) error {
	// Get the MQTT Control Packet type.
	ptype, err := p.Type()
	if err != nil {
		return err
	}

	switch ptype {
	case packet.TypeCONNACK:
		// Notify the arrival of the CONNACK Packet if possible.
		select {
		case cli.connackc <- struct{}{}:
		default:
		}
	case
		packet.TypePUBLISH,
		packet.TypePUBACK,
		packet.TypePUBREC,
		packet.TypePUBREL,
		packet.TypePUBCOMP,
		packet.TypeSUBACK,
		packet.TypeUNSUBACK,
		packet.TypePINGRESP:
	default:
		return packet.ErrInvalidPacketType
	}

	return nil
}

// handleError handles the error and disconnects
// the Network Connection.
func (cli *Client) handleErrorAndDisconn(err error) {
	// Lock for reading.
	cli.mu.RLock()

	// Ignore the error and end the process
	// if the Network Connection has already
	// disconnected.
	if cli.conn == nil || cli.conn.disconnected {
		// Unlock.
		cli.mu.RUnlock()

		return
	}

	// Unlock.
	cli.mu.RUnlock()

	// Handle the error.
	if cli.errHandler != nil {
		cli.errHandler(err)
	}

	// Send a disconnect signal to the goroutine
	// via the channel if possible.
	select {
	case cli.disconnc <- struct{}{}:
	default:
	}
}

// sendPackets sends Packets to the Server.
func (cli *Client) sendPackets(keepAlive uint16) {
	defer cli.wgConn.Done()

	for {
		var keepAlivec <-chan time.Time

		if keepAlive > 0 {
			keepAlivec = time.After(time.Duration(keepAlive) * time.Second)
		}

		select {
		case p := <-cli.sendc:
			cli.sendPacket(p)
		case <-keepAlivec:
			cli.sendPacket(packet.NewPINGREQ())
		case <-cli.sendEndc:
			return
		}
	}
}

// sendPacket sends a Packet to the Server.
func (cli *Client) sendPacket(p packet.Packet) {
	// Lock for sending the Packet.
	cli.mu.RLock()

	// Send the Packet to the Server.
	err := cli.send(p)

	// Unlock.
	cli.mu.RUnlock()

	if err != nil {
		cli.handleErrorAndDisconn(err)
	}
}

// New creates and returns a Client.
func New(opts *Options) *Client {
	// Initialize the options.
	if opts == nil {
		opts = &Options{}
	}
	// Create a Client.
	cli := &Client{
		disconnc:    make(chan struct{}, 1),
		disconnEndc: make(chan struct{}, 1),
		errHandler:  opts.ErrHandler,
	}

	// Initialize the channels of the client.
	cli.initChans()

	// Launch a goroutine which disconnects the Network Connection.
	cli.wgNew.Add(1)
	go func() {
		defer func() {
			cli.wgNew.Done()
		}()

		for {
			select {
			case <-cli.disconnc:
			case <-cli.disconnEndc:
				// End the goroutine.
				return
			}
		}

	}()

	// Return the Client.
	return cli
}
