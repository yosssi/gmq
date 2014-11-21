package cmd

import (
	"fmt"

	"github.com/yosssi/gmq/common"
)

// Defalut values of the flags
const (
	defConnHost         = "localhost"
	defConnPort         = 1883
	defConnCleanSession = true
	defConnWill         = false
	defConnWillQoS      = common.QoS0
	defConnWillRetain   = false
	defConnUsername     = ""
	defConnPassword     = ""
)

// Conn is a command which sends a connection request to the server.
var Conn = &Cmd{
	Name:  "conn",
	Usage: "send a connection request to the server",
	Run:   conn,
}

// Flags of the Conn command
var (
	connHost         string
	connPort         uint
	connCleanSession bool
	connWill         bool
	connWillQoS      uint
	connWillRetain   bool
	connUsername     string
	connPassword     string
)

func init() {
	Conn.Flag.StringVar(&connHost, "h", defConnHost, "host name of the server to connect to")
	Conn.Flag.UintVar(&connPort, "p", defConnPort, "port number of the server to connect to")
	Conn.Flag.BoolVar(&connCleanSession, "c", defConnCleanSession, "enable crean session")
	Conn.Flag.BoolVar(&connWill, "w", defConnWill, "enable Will")
	Conn.Flag.UintVar(&connWillQoS, "wq", defConnWillQoS, "QoS levels to be used when publishing the Will Message")
	Conn.Flag.BoolVar(&connWillRetain, "wr", defConnWillRetain, "enable Will Retain")
	Conn.Flag.StringVar(&connUsername, "u", defConnUsername, "username")
	Conn.Flag.StringVar(&connPassword, "P", defConnPassword, "password")
}

// conn sends a connection request to the server.
func conn(c *Cmd) error {
	fmt.Println(connHost, connPort, connCleanSession, connWill, connWillQoS, connWillRetain, connUsername, connPassword)
	return nil
}
