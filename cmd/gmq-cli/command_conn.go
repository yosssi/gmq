package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Default values
const (
	defaultNetwork             = "tcp"
	defaultHost                = "localhost"
	defaultPort           uint = 1883
	defaultCONNACKTimeout uint = 30
)

// Timeout in seconds for sending the connacked signal
const connackedSigTimeout = 1

// Hostname
var hostname, _ = os.Hostname()

// Error value
var errCONNACKTimeout = errors.New("the CONNACK Packet was not received within a reasonalbe amount of time")

// commandConn represents a conn command.
type commandConn struct {
	ctx            *context
	network        string
	address        string
	connacked      chan struct{}
	connackTimeout time.Duration
	connectOpts    *packet.CONNECTOptions
}

// run tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cmd *commandConn) run() error {
	// Lock for the connection.
	cmd.ctx.mu.Lock()

	// Try  to establish a Network Connection to the Server and
	// send a CONNECT Packet to the Server.
	err := cmd.ctx.cli.Connect(cmd.network, cmd.address, cmd.connectOpts)

	// Unlock.
	cmd.ctx.mu.Unlock()

	if err != nil {
		return err
	}

	// Launch a goroutine which waits for receiving the CONNACK Packet.
	cmd.ctx.wg.Add(1)
	go cmd.waitCONNACK()

	// Launch a goroutine which receives a Packet from the Server.
	cmd.ctx.wg.Add(1)
	go cmd.receive()

	// Launch a goroutine which sends a Packet to the Server.
	cmd.ctx.wg.Add(1)
	go cmd.send()

	return nil
}

// waitCONNACK waits for receiving the CONNACK Packet.
func (cmd *commandConn) waitCONNACK() {
	defer cmd.ctx.wg.Done()

	var timeout <-chan time.Time

	if cmd.connackTimeout > 0 {
		timeout = time.After(cmd.connackTimeout * time.Second)
	}

	select {
	case <-cmd.ctx.connack:
	case <-timeout:
		printError(errCONNACKTimeout)

		// Send a disconnect signal to the channel if possible.
		cmd.ctx.disconn <- struct{}{}
	case <-cmd.ctx.connackEnd:
	}
}

// receive receives a Packet from the Server.
func (cmd *commandConn) receive() {
	defer cmd.ctx.wg.Done()

	for {
		// Receive a Packet from the Network Connection.
		p, err := cmd.ctx.cli.Receive()
		if err != nil {
			printError(err)

			// Send a disconnect signal to the channel if possible.
			select {
			case cmd.ctx.disconn <- struct{}{}:
			default:
			}

			return
		}

		// Handle the Packet.
		if err := cmd.handle(p); err != nil {
			printError(err)
		}
	}
}

// handle handles the Packet.
func (cmd *commandConn) handle(p packet.Packet) error {
	// Get the MQTT Control Packet type.
	ptype, err := p.Type()
	if err != nil {
		return err
	}

	switch ptype {
	case packet.TypeCONNACK:
		// Notify the arrival of the CONNACK Packet if possible.
		select {
		case cmd.ctx.connack <- struct{}{}:
		default:
		}
	}

	return nil
}

// send sends the Packet to the Server.
func (cmd *commandConn) send() {
	defer cmd.ctx.wg.Done()

	for {
		var keepAlive <-chan time.Time

		if *cmd.connectOpts.KeepAlive > 0 {
			keepAlive = time.After(time.Duration(*cmd.connectOpts.KeepAlive) * time.Second)
		}

		select {
		case p := <-cmd.ctx.send:
			cmd.sendPacket(p)
		case <-keepAlive:
			cmd.sendPacket(packet.NewPINGREQ())
		case <-cmd.ctx.sendEnd:
			return
		}
	}
}

// sendPacket sends the Packet to the Server.
func (cmd *commandConn) sendPacket(p packet.Packet) {
	// Lock for sending the Packet.
	cmd.ctx.mu.RLock()

	// Send the Packet to the Server.
	err := cmd.ctx.cli.Send(p)

	// Unlock.
	cmd.ctx.mu.RUnlock()

	if err != nil {
		printError(err)
	}
}

// newCommandConn creates and returns a conn command.
func newCommandConn(args []string, ctx *context) (*commandConn, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	network := flg.String("n", defaultNetwork, "network on which the Client connects to the Server")
	host := flg.String("h", defaultHost, "host name of the Server which the Client connects to")
	port := flg.Uint("p", defaultPort, "port number of the Server which the Client connects to")
	connackTimeout := flg.Uint(
		"ackt",
		defaultCONNACKTimeout,
		"Timeout in seconds for the Client to wait for receiving the CONNACK Packet after sending the CONNECT Packet",
	)
	clientID := flg.String("i", hostname, "Client identifier for the Client")
	cleanSession := flg.Bool("c", packet.DefaultCleanSession, "Clean Session")
	willTopic := flg.String("wt", "", "Will Topic")
	willMessage := flg.String("wm", "", "Will Message")
	willQoS := flg.Uint("wq", mqtt.QoS0, "Will QoS")
	willRetain := flg.Bool("wr", false, "Will Retain")
	userName := flg.String("u", "", "User Name")
	password := flg.String("P", "", "Password")
	keepAlive := flg.Uint("k", packet.DefaultKeepAlive, "Keep Alive measured in seconds")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	// Create a conn command.
	cmd := &commandConn{
		ctx:            ctx,
		network:        *network,
		address:        *host + ":" + strconv.Itoa(int(*port)),
		connacked:      make(chan struct{}),
		connackTimeout: time.Duration(*connackTimeout),
		connectOpts: &packet.CONNECTOptions{
			ClientID:     *clientID,
			CleanSession: cleanSession,
			WillTopic:    *willTopic,
			WillMessage:  *willMessage,
			WillQoS:      *willQoS,
			WillRetain:   *willRetain,
			UserName:     *userName,
			Password:     *password,
			KeepAlive:    keepAlive,
		},
	}

	return cmd, nil
}
