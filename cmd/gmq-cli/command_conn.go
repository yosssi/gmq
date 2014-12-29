package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Command name
const cmdNameConn = "conn"

// Default values
const (
	defaultNetwork              = "tcp"
	defaultHost                 = "localhost"
	defaultPort            uint = 1883
	defaultCONNACKTimeout  uint = 30
	defaultPINGRESPTimeout uint = 30
	defaultCleanSession         = true
	defaultKeepAlive       uint = 60
)

// Error value
var errParseCrtFailure = errors.New("failed to parse root certificate")

// commandConn represents a conn command.
type commandConn struct {
	cli         *client.Client
	connectOpts *client.ConnectOptions
}

// run tries to establish a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cmd *commandConn) run() error {
	return cmd.cli.Connect(cmd.connectOpts)
}

// newCommandConn creates and returns a conn command.
func newCommandConn(args []string, cli *client.Client) (command, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	network := flg.String("n", defaultNetwork, "network on which the Client connects to the Server")
	host := flg.String("h", defaultHost, "host name of the Server which the Client connects to")
	port := flg.Uint("p", defaultPort, "port number of the Server which the Client connects to")
	crtPath := flg.String("crt", "", "the path of the certificate authority file to verify the server connection")
	connackTimeout := flg.Uint(
		"ct",
		defaultCONNACKTimeout,
		"Timeout in seconds for the Client to wait for receiving the CONNACK Packet after sending the CONNECT Packet",
	)
	pingrespTimeout := flg.Uint(
		"pt",
		defaultPINGRESPTimeout,
		"Timeout in seconds for the Client to wait for receiving the PINGRESP Packet after sending the PINGREQ Packet",
	)
	clientID := flg.String("i", "", "Client identifier for the Client")
	userName := flg.String("u", "", "User Name")
	password := flg.String("P", "", "Password")
	cleanSession := flg.Bool("c", defaultCleanSession, "Clean Session")
	keepAlive := flg.Uint("k", defaultKeepAlive, "Keep Alive measured in seconds")
	willTopic := flg.String("wt", "", "Will Topic")
	willMessage := flg.String("wm", "", "Will Message")
	willQoS := flg.Uint("wq", uint(mqtt.QoS0), "Will QoS")
	willRetain := flg.Bool("wr", false, "Will Retain")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	var tlsConfig *tls.Config

	// Parse the certificate authority file.
	if *crtPath != "" {
		b, err := ioutil.ReadFile(*crtPath)
		if err != nil {
			return nil, err
		}

		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(b); !ok {
			return nil, errParseCrtFailure
		}

		tlsConfig = &tls.Config{
			RootCAs: roots,
		}
	}

	// Create a conn command.
	cmd := &commandConn{
		cli: cli,
		connectOpts: &client.ConnectOptions{
			Network:         *network,
			Address:         *host + ":" + strconv.Itoa(int(*port)),
			TLSConfig:       tlsConfig,
			CONNACKTimeout:  time.Duration(*connackTimeout),
			PINGRESPTimeout: time.Duration(*pingrespTimeout),
			ClientID:        []byte(*clientID),
			UserName:        []byte(*userName),
			Password:        []byte(*password),
			CleanSession:    *cleanSession,
			KeepAlive:       uint16(*keepAlive),
			WillTopic:       []byte(*willTopic),
			WillMessage:     []byte(*willMessage),
			WillQoS:         byte(*willQoS),
			WillRetain:      *willRetain,
		},
	}

	// Return the command.
	return cmd, nil
}
