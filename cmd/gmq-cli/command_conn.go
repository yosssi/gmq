package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Default values
const (
	defaultHost      = "localhost"
	defaultPort uint = 1883
)

// Hostname
var hostname, _ = os.Hostname()

// commandConn represents a conn command.
type commandConn struct {
	ctx         *context
	opts        *client.ConnectOptions
	connectOpts *packet.CONNECTOptions
}

// run tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cmd *commandConn) run() error {
	// Try to establish a Network Connection to the Server.
	if err := cmd.ctx.cli.Connect(cmd.opts); err != nil {
		return err
	}

	// Send a CONNECT Packet to the Server.
	if err := cmd.ctx.cli.SendCONNECT(cmd.connectOpts); err != nil {
		// Disconnect the Network Connection.
		if anotherErr := cmd.ctx.cli.Disconnect(); anotherErr != nil {
			return fmt.Errorf(strErrMulti, anotherErr, err)
		}

		return err
	}

	// Launch a goroutine which sends a Packet to the Server.

	// Launch a goroutine which receives a Packet from the Server.

	// Launch a goroutine which reads data from the Network Connection.
	go func() {
		for {
			ptype, p, err := cmd.ctx.cli.Receive()
			fmt.Println(ptype, p, err)
		}
	}()

	return nil
}

// newCommandConn creates and returns a conn command.
func newCommandConn(args []string, ctx *context) (*commandConn, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	network := flg.String("n", client.DefaultNetwork, "network on which the Client connects to the Server")
	host := flg.String("h", defaultHost, "host name of the Server to connect to")
	port := flg.Uint("p", defaultPort, "port number of the Server to connect to")
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
		ctx: ctx,
		opts: &client.ConnectOptions{
			Network: *network,
			Address: *host + ":" + strconv.Itoa(int(*port)),
		},
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
