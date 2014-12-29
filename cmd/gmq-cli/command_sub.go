package main

import (
	"flag"
	"os"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Command name
const cmdNameSub = "sub"

// commandSub represents a sub command.
type commandSub struct {
	cli           *client.Client
	subscribeOpts *client.SubscribeOptions
}

// run sends a SUBSCRIBE Packet to the Server.
func (cmd *commandSub) run() error {
	return cmd.cli.Subscribe(cmd.subscribeOpts)
}

// newCommandSub creates and returns a sub command.
func newCommandSub(args []string, cli *client.Client) (command, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	topicFilter := flg.String("t", "", "Topic Filter")
	qos := flg.Uint("q", uint(mqtt.QoS0), "QoS")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	// Create a sub command.
	cmd := &commandSub{
		cli: cli,
		subscribeOpts: &client.SubscribeOptions{
			SubReqs: []*client.SubReq{
				&client.SubReq{
					TopicFilter: []byte(*topicFilter),
					QoS:         byte(*qos),
					Handler:     messageHandler,
				},
			},
		},
	}

	// Return the command.
	return cmd, nil
}

func messageHandler(topicName, message []byte) {
	os.Stdout.WriteString("\n[Topic Name]\n" + string(topicName) + "\n[Application Message]\n" + string(message) + "\n")
	printHeader()
}
