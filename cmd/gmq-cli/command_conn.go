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
	cmd.ctx.climu.Lock()
	defer cmd.ctx.climu.Unlock()

	// Try to establish a Network Connection to the Server and
	// send a CONNECT Packet to the Server.
	if err := cmd.ctx.cli.Connect(cmd.network, cmd.address, cmd.connectOpts); err != nil {
		return err
	}

	// Launch a goroutine which sends a Packet to the Server.
	go func() {
		defer func() {
			// Send an ended signal.
			cmd.ctx.sendEndedc <- struct{}{}
		}()

	SendLoop:
		for {
			var keepAlive <-chan time.Time
			if *cmd.connectOpts.KeepAlive > 0 {
				keepAlive = time.After(time.Duration(*cmd.connectOpts.KeepAlive) * time.Second)
			}

			select {
			case p := <-cmd.ctx.sendc:
				// Send the Packet to the Server.
				if err := sendWithLock(cmd.ctx, p); err != nil {
					// Disconnect the Network Connection.
					go func() {
						if err := disconnectWithLock(cmd.ctx); err != nil {
							cmd.ctx.errc <- err
							return
						}
					}()

					cmd.ctx.errc <- err

					break SendLoop
				}
			case <-keepAlive:
				// Send a PINGREQ Packet to the Server.
				if err := sendWithLock(cmd.ctx, packet.NewPINGREQ()); err != nil {
					// Disconnect the Network Connection.
					go func() {
						if err := disconnectWithLock(cmd.ctx); err != nil {
							cmd.ctx.errc <- err
							return
						}
					}()

					cmd.ctx.errc <- err

					break SendLoop
				}
			case <-cmd.ctx.sendEndc:
				return
			}
		}

		<-cmd.ctx.sendEndc
	}()

	// Launch a goroutine which reads data from the Network Connection.
	go func() {
		defer func() {
			// Send an ended signal.
			cmd.ctx.readEndedc <- struct{}{}
		}()

		for {
			p, err := cmd.ctx.cli.Receive()
			if err != nil {
				// Disconnect the Network Connection.
				go func() {
					if err := disconnectWithLock(cmd.ctx); err != nil {
						cmd.ctx.errc <- err
						return
					}
				}()

				cmd.ctx.errc <- err

				return
			}

			cmd.ctx.recvc <- p
		}
	}()

	// Launch a goroutine which monitors the arrival of the CONNACK Packet if connackTimeout > 0.
	go func() {
		defer func() {
			// Send an ended signal.
			cmd.ctx.connackEndedc <- struct{}{}
		}()

		var timeout <-chan time.Time

		if cmd.connackTimeout > 0 {
			timeout = time.After(cmd.connackTimeout * time.Second)
		}

		select {
		case <-cmd.ctx.connackc:
		case <-timeout:
			// Disconnect the Network Connection.
			go func() {
				if err := disconnectWithLock(cmd.ctx); err != nil {
					cmd.ctx.errc <- err
					return
				}
			}()

			cmd.ctx.errc <- errCONNACKTimeout
		case <-cmd.ctx.connackEndc:
			return
		}

		<-cmd.ctx.connackEndc
	}()

	// Launch a goroutine which handles a Packet received from the Server.
	go func() {
		defer func() {
			// Send an ended signal.
			cmd.ctx.recvEndedc <- struct{}{}
		}()

		for {
			select {
			case p := <-cmd.ctx.recvc:
				ptype, err := p.Type()
				if err != nil {
					cmd.ctx.errc <- err
					continue
				}

				switch ptype {
				case packet.TypeCONNACK:
					// Notify the arrival of the CONNACK Packet.
					//cmd.ctx.connackc <- struct{}{}
				}
			case <-cmd.ctx.recvEndc:
				return
			}
		}
	}()

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
