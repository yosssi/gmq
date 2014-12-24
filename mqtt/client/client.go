package client

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Multiple errors string format
const strErrMulti = "error (%q) occurred while handling the other error (%q)"

// Maximum Packet Identifier
const maxPacketID = 65535

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
	ErrCONNACKTimeout   = errors.New("the CONNACK Packet was not received within a reasonalbe amount of time")
	ErrPINGRESPTimeout  = errors.New("the PINGRESP Packet was not received within a reasonalbe amount of time")
	ErrPacketIDExhaused = errors.New("Packet Identifiers are exhausted")
)

// Client represents a Client.
type Client struct {
	// muConn is the Mutex for the Network Connection.
	muConn sync.RWMutex
	// conn is the Network Connection.
	conn *connection

	// muSess is the Mutex for the Session.
	muSess sync.RWMutex
	// sess is the Session.
	sess *session

	// wg is the Wait Group for the goroutines
	// which are launched by the New method.
	wg sync.WaitGroup
	// disconnc is the channel which handles the signal
	// to disconnect the Network Connection.
	disconnc chan struct{}
	// disconnEndc is the channel which ends the goroutine
	// which disconnects the Network Connection.
	disconnEndc chan struct{}

	// errHandler is the error handler.
	errHandler func(error)
}

// Connect establishes a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(opts *ConnectOptions) error {
	// Lock for the connection.
	cli.muConn.Lock()

	// Unlock.
	defer cli.muConn.Unlock()

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

	// Lock for reading and updating the Session.
	cli.muSess.Lock()

	// Create a Session or reuse the current Session.
	if opts.CleanSession || cli.sess == nil {
		// Create a Session and set it to the Client.
		cli.sess = newSession(opts.CleanSession, opts.ClientID)
	} else {
		// Reuse the Session and set its Client Identifier to the options.
		opts.ClientID = cli.sess.clientID
	}

	// Unlock.
	cli.muSess.Unlock()

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
	cli.conn.wg.Add(1)
	go cli.waitPacket(cli.conn.connack, opts.CONNACKTimeout, ErrCONNACKTimeout)

	// Launch a goroutine which receives a Packet from the Server.
	cli.conn.wg.Add(1)
	go cli.receivePackets()

	// Launch a goroutine which sends a Packet to the Server.
	cli.conn.wg.Add(1)
	go cli.sendPackets(time.Duration(opts.KeepAlive), opts.PINGRESPTimeout)

	return nil
}

// Disconnect sends a DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Lock for the disconnection.
	cli.muConn.Lock()

	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		// Unlock.
		cli.muConn.Unlock()

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
		cli.muConn.Unlock()

		return err
	}

	// Change the state of the Network Connection to disconnected.
	cli.conn.disconnected = true

	// Send the end signal to the goroutine via the channels.
	select {
	case cli.conn.sendEnd <- struct{}{}:
	default:
	}

	// Unlock.
	cli.muConn.Unlock()

	// Wait until all goroutines end.
	cli.conn.wg.Wait()

	// Lock for cleaning the Network Connection.
	cli.muConn.Lock()

	// Lock for cleaning the Session.
	cli.muSess.Lock()

	// Clean the Network Connection and the Session.
	cli.clean()

	// Unlock.
	cli.muSess.Unlock()

	// Unlock.
	cli.muConn.Unlock()

	return nil
}

// Publish sends a PUBLISH Packet to the Server.
func (cli *Client) Publish(opts *PublishOptions) error {
	// Lock for reading.
	cli.muConn.RLock()

	// Unlock.
	defer cli.muConn.RUnlock()

	// Check the Network Connection.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &PublishOptions{}
	}

	// Define a Packet Identifier.
	var packetID uint16

	if opts.QoS != mqtt.QoS0 {
		// Lock for reading and updating the Session.
		cli.muSess.Lock()

		// Define an error.
		var err error

		// Generate a Packet Identifer.
		packetID, err = cli.generatePacketID()
		if err != nil {
			// Unlock.
			cli.muSess.Unlock()

			return err
		}
	}

	// Create a PUBLISH Packet.
	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:       opts.QoS,
		Retain:    opts.Retain,
		TopicName: []byte(opts.TopicName),
		PacketID:  packetID,
		Message:   []byte(opts.Message),
	})
	if err != nil {
		if opts.QoS != mqtt.QoS0 {
			// Unlock.
			cli.muSess.Unlock()
		}

		return err
	}

	if opts.QoS != mqtt.QoS0 {
		// Set the Packet to the Session.
		cli.sess.sendingPackets[packetID] = p

		// Unlock.
		cli.muSess.Unlock()
	}

	// Send the Packet to the Server.
	cli.conn.send <- p

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

// waitPacket waits for receiving the Packet.
func (cli *Client) waitPacket(packetc <-chan struct{}, timeout time.Duration, errTimeout error) {
	defer cli.conn.wg.Done()

	var timeoutc <-chan time.Time

	if timeout > 0 {
		timeoutc = time.After(timeout * time.Second)
	}

	select {
	case <-packetc:
	case <-timeoutc:
		// Handle the timeout error.
		cli.handleErrorAndDisconn(errTimeout)

		// Sned a disconnect siganl to the goroutine
		// via the channel if possible.
		select {
		case cli.disconnc <- struct{}{}:
		default:
		}
	}
}

// receivePackets receives Packets from the Server.
func (cli *Client) receivePackets() {
	defer func() {
		// Close the channel which handles a signal which
		// notifies the arrival of the CONNACK Packet.
		close(cli.conn.connack)

		cli.conn.wg.Done()
	}()

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
			// Handle the error and disconnect
			// the Network Connection.
			cli.handleErrorAndDisconn(err)

			// End the goroutine.
			return
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
		case cli.conn.connack <- struct{}{}:
		default:
		}
	case packet.TypePINGRESP:
		// Lock for reading and updating pingrespcs.
		cli.conn.muPINGRESPs.Lock()

		// Check the length of pingrespcs.
		if len(cli.conn.pingresps) == 0 {
			// End the function if there is no channel in pingrespcs.
			return nil
		}

		// Get the first channel in pingrespcs.
		pingrespc := cli.conn.pingresps[0]

		// Remove the first channel from pingrespcs.
		cli.conn.pingresps = cli.conn.pingresps[1:]

		// Unlock.
		cli.conn.muPINGRESPs.Unlock()

		// Notify the arrival of the PINGRESP Packet if possible.
		select {
		case pingrespc <- struct{}{}:
		default:
		}
	case
		packet.TypePUBLISH,
		packet.TypePUBACK,
		packet.TypePUBREC,
		packet.TypePUBREL,
		packet.TypePUBCOMP,
		packet.TypeSUBACK,
		packet.TypeUNSUBACK:
	default:
		return packet.ErrInvalidPacketType
	}

	return nil
}

// handleError handles the error and disconnects
// the Network Connection.
func (cli *Client) handleErrorAndDisconn(err error) {
	// Lock for reading.
	cli.muConn.RLock()

	// Ignore the error and end the process
	// if the Network Connection has already
	// been disconnected.
	if cli.conn == nil || cli.conn.disconnected {
		// Unlock.
		cli.muConn.RUnlock()

		return
	}

	// Unlock.
	cli.muConn.RUnlock()

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
func (cli *Client) sendPackets(keepAlive time.Duration, pingrespTimeout time.Duration) {
	defer func() {
		// Lock for reading and updating pingrespcs.
		cli.conn.muPINGRESPs.Lock()

		// Close the channels which handle a signal which
		// notifies the arrival of the PINGREQ Packet.
		for _, pingresp := range cli.conn.pingresps {
			close(pingresp)
		}

		// Initialize pingrespcs
		cli.conn.pingresps = make([]chan struct{}, 0)

		// Unlock.
		cli.conn.muPINGRESPs.Unlock()

		cli.conn.wg.Done()
	}()

	for {
		var keepAlivec <-chan time.Time

		if keepAlive > 0 {
			keepAlivec = time.After(keepAlive * time.Second)
		}

		select {
		case p := <-cli.conn.send:
			// Lock for sending the Packet.
			cli.muConn.RLock()

			// Send the Packet to the Server.
			err := cli.send(p)

			// Unlock.
			cli.muConn.RUnlock()

			if err != nil {
				// Handle the error and disconnect the Network Connection.
				cli.handleErrorAndDisconn(err)

				// End this function.
				return
			}
		case <-keepAlivec:
			// Lock for sending the Packet.
			cli.muConn.RLock()

			// Send a PINGREQ Packet to the Server.
			err := cli.send(packet.NewPINGREQ())

			// Unlock.
			cli.muConn.RUnlock()

			if err != nil {
				// Handle the error and disconnect the Network Connection.
				cli.handleErrorAndDisconn(err)

				// End this function.
				return
			}

			// Create a channel which handles the signal to notify the arrival of
			// the PINGRESP Packet.
			pingresp := make(chan struct{})

			// Lock for appending the channel to pingrespcs.
			cli.conn.muPINGRESPs.Lock()

			// Append the channel to pingrespcs.
			cli.conn.pingresps = append(cli.conn.pingresps, pingresp)

			// Unlock.
			cli.conn.muPINGRESPs.Unlock()

			// Launch a goroutine which waits for receiving the PINGRESP Packet.
			cli.conn.wg.Add(1)
			go cli.waitPacket(pingresp, pingrespTimeout, ErrPINGRESPTimeout)
		case <-cli.conn.sendEnd:
			// End this function.
			return
		}
	}
}

// generatePacketID generates and returns a Packet Identifier.
func (cli *Client) generatePacketID() (uint16, error) {
	// Define a Packet Identifier.
	var id uint16

	for {
		// Find a Packet Identifier which does not used.
		if _, exist := cli.sess.sendingPackets[id]; !exist {
			// Return the Packet Identifier.
			return id, nil
		}

		if id == maxPacketID {
			break
		}

		id++
	}

	// Return an error if available ids are not found.
	return 0, ErrPacketIDExhaused
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

	// Launch a goroutine which disconnects the Network Connection.
	cli.wg.Add(1)
	go func() {
		defer func() {
			cli.wg.Done()
		}()

		for {
			select {
			case <-cli.disconnc:
				if err := cli.Disconnect(); err != nil {
					if cli.errHandler != nil {
						cli.errHandler(err)
					}
				}
			case <-cli.disconnEndc:
				// End the goroutine.
				return
			}
		}

	}()

	// Return the Client.
	return cli
}
