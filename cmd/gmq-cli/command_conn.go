package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
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

// Error value
var errCONNACKTimeout = errors.New("the Network Connection was disconnected because the CONNACK Packet was not received within a reasonalbe amount of time")

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
	// Get a lock for the Network Connection.
	cmd.ctx.connMu.Lock()
	defer cmd.ctx.connMu.Unlock()

	// Check the state of the Network Connection.
	if atomic.LoadUint32(&cmd.ctx.connState) != connStateClosed {
		return client.ErrAlreadyConnected
	}

	// Try to establish a Network Connection to the Server and
	// send a CONNECT Packet to the Server.
	if err := cmd.ctx.cli.Connect(cmd.network, cmd.address, cmd.connectOpts); err != nil {
		return err
	}

	return nil
}

// newCommandConn creates and returns a conn command.
func newCommandConn(args []string, ctx *context) (*commandConn, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	network := flg.String("n", defaultNetwork, "network on which the Client connects to the Server")
	host := flg.String("h", defaultHost, "host name of the Server to connect to")
	port := flg.Uint("p", defaultPort, "port number of the Server to connect to")
	connackTimeout := flg.Uint(
		"ackt",
		defaultCONNACKTimeout,
		"Timeout in seconds for the Client to wait receiving the CONNACK Packet after sending the CONNECT Packet",
	)
	clientID := flg.String("i", hostname, "Client identifier for the Client")
	cleanSession := flg.Bool("c", packet.DefaultCleanSession, "Clean Session")
	willTopic := flg.String("wt", "", "Will Topic")
	willMessage := flg.String("wm", "", "Will Message")
	willQoS := flg.Uint("wq", mqtt.QoS0, "Will QoS")
	willRetain := flg.Bool("wr", false, "Will Retain")
	userName := flg.String("u", "", "User Name")
	password := flg.String("P", "", "Password")
	keepAlive := flg.Uint("k", packet.DefaultKeepAlive, "Keep Alive in seconds for the Client")

	// Parse the flag definitions from the arguments.
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

	// Return the command.
	return cmd, nil
}
