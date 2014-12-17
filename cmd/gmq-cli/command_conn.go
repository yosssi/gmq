package main

import (
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
	return nil
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
