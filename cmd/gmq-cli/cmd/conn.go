package cmd

import (
	"os"
	"strconv"

	"github.com/yosssi/gmq/client"
	"github.com/yosssi/gmq/common"
	"github.com/yosssi/gmq/common/packet"
)

// Conn is a command which sends a connection request to the Server.
var Conn = &Cmd{
	Name:  "conn",
	Usage: "send a connection request to the Server",
	Run:   conn,
}

// Default values
const (
	defaultHost      = "localhost"
	defaultPort uint = 1883
)

// Flags of the Conn command
var (
	connHost         string
	connPort         uint
	connClientID     string
	connCleanSession bool
	connWillTopic    string
	connWillMessage  string
	connWillQoS      uint
	connWillRetain   bool
	connUserName     string
	connPassword     string
	connKeepAlive    uint
)

func init() {
	hostname, _ := os.Hostname()

	Conn.Flag.StringVar(&connHost, "h", defaultHost, "host name of the server to connect to")
	Conn.Flag.UintVar(&connPort, "p", uint(defaultPort), "port number of the server to connect to")
	Conn.Flag.StringVar(&connClientID, "i", hostname, "id to use for this client")
	Conn.Flag.BoolVar(&connCleanSession, "c", packet.DefaultCleanSession, "Clean Session")
	Conn.Flag.StringVar(&connWillTopic, "wc", "", "Will Topic")
	Conn.Flag.StringVar(&connWillMessage, "wm", "", "Will Message")
	Conn.Flag.UintVar(&connWillQoS, "wq", common.QoS0, "Will QoS")
	Conn.Flag.BoolVar(&connWillRetain, "wr", false, "Will Retain")
	Conn.Flag.StringVar(&connUserName, "u", "", "User Name")
	Conn.Flag.StringVar(&connPassword, "P", "", "Password")
	Conn.Flag.UintVar(&connKeepAlive, "k", uint(packet.DefaultKeepAlive), "Keep Alive in seconds for this client.")
}

// connect sends a connection request to the server.
func conn(cli *client.Client, c *Cmd) error {
	return cli.Connect(
		&client.ConnectOptions{
			Network:        client.DefaultNetwork,
			Address:        connHost + ":" + strconv.Itoa(int(connPort)),
			CONNACKTimeout: client.DefaultCONNACKTimeout,
		},
		&packet.CONNECTOptions{
			ClientID:     connClientID,
			CleanSession: &connCleanSession,
			WillTopic:    connWillTopic,
			WillMessage:  connWillMessage,
			WillQoS:      connWillQoS,
			WillRetain:   connWillRetain,
			UserName:     connUserName,
			Password:     connPassword,
			KeepAlive:    &connKeepAlive,
		})
}
