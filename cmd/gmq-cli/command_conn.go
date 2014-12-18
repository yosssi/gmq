package main

import (
	"flag"
	"fmt"
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

// Hostname
var hostname, _ = os.Hostname()

// commandConn represents a conn command.
type commandConn struct {
	ctx            *context
	network        string
	address        string
	connackTimeout time.Duration
	connectOpts    *packet.CONNECTOptions
}

// run tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cmd *commandConn) run() error {
	// Try to establish a Network Connection to the Server and
	// send a CONNECT Packet to the Server.
	if err := cmd.connect(); err != nil {
		return err
	}

	// Launch a goroutine which receives an MQTT Control Packet
	// from the Server.
	cmd.ctx.wg.Add(1)
	go cmd.receive(cmd.ctx.wg.Done)

	return nil
}

// connect tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cmd *commandConn) connect() error {
	// Lock for connecting to the Server.
	cmd.ctx.mu.Lock()
	defer cmd.ctx.mu.Unlock()

	return cmd.ctx.cli.Connect(cmd.network, cmd.address, cmd.connectOpts)
}

// receive receives an MQTT Control Packet from the Server.
func (cmd *commandConn) receive(f func()) {
	if f != nil {
		defer f()
	}

	for {
		// Receive an MQTT Control Packet from the Server.
		p, err := cmd.ctx.cli.Receive()
		if err != nil {
			// Ignore the error and end this function if the Network Connection is not connected.
			if !cmd.ctx.getConnected() {
				return
			}

			// Print the error.
			printError(err)

			// Disconnect the Network Connection.
			if err := disconnect(cmd.ctx); err != nil {
				printError(err)
			}

			return
		}

		// Launch a goroutine which handles the Packet.
		cmd.ctx.wg.Add(1)
		go cmd.handle(p, cmd.ctx.wg.Done)
	}
}

// handle handles an MQTT Control Packet.
func (cmd *commandConn) handle(p packet.Packet, f func()) {
	if f != nil {
		defer f()
	}

	// Get the Packet type.
	ptype, err := p.Type()
	if err != nil {
		printError(err)
		return
	}

	switch ptype {
	case packet.TypeCONNACK:
		// TODO
		fmt.Println("CONNACK!!!")
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
