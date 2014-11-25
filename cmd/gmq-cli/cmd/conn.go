package cmd

import (
	"strconv"

	"github.com/yosssi/gmq/client"
	"github.com/yosssi/gmq/common"
)

// Conn is a command which sends a connection request to the server.
var Conn = &Cmd{
	Name:  "conn",
	Usage: "send a connection request to the server",
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
	Conn.Flag.StringVar(&connHost, "h", defaultHost, "host name of the server to connect to")
	Conn.Flag.UintVar(&connPort, "p", uint(defaultPort), "port number of the server to connect to")
	Conn.Flag.BoolVar(&connCleanSession, "c", common.DefaultCleanSession, "Clean Session")
	Conn.Flag.StringVar(&connWillTopic, "wc", "", "Will Topic")
	Conn.Flag.StringVar(&connWillMessage, "wm", "", "Will Message")
	Conn.Flag.UintVar(&connWillQoS, "wq", common.QoS0, "Will QoS")
	Conn.Flag.BoolVar(&connWillRetain, "wr", false, "Will Retain")
	Conn.Flag.StringVar(&connUserName, "u", "", "User Name")
	Conn.Flag.StringVar(&connPassword, "P", "", "Password")
	Conn.Flag.UintVar(&connKeepAlive, "k", uint(common.DefaultKeepAlive), "Keep Alive in seconds for this client.")
}

// connect sends a connection request to the server.
func conn(cli *client.Client, c *Cmd) error {
	return cli.Connect(
		connHost+":"+strconv.Itoa(int(connPort)),
		&common.OptionsPacketCONNECT{
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
