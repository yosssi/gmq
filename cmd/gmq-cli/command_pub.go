package main

import (
	"flag"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// Command name
const cmdNamePub = "pub"

// commandPub represents a pub command.
type commandPub struct {
	cli         *client.Client
	publishOpts *client.PublishOptions
}

// Publish sends a PUBLISH Packet to the Server.
func (cmd *commandPub) run() error {
	return cmd.cli.Publish(cmd.publishOpts)
}

// newCommandPub creates and returns a pub command.
func newCommandPub(args []string, cli *client.Client) (command, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	qos := flg.Uint("q", uint(mqtt.QoS0), "QoS")
	retain := flg.Bool("r", false, "Retain")
	topicName := flg.String("t", "", "Topic Name")
	message := flg.String("m", "", "Application Message")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	// Create a pub command.
	cmd := &commandPub{
		cli: cli,
		publishOpts: &client.PublishOptions{
			QoS:       byte(*qos),
			Retain:    *retain,
			TopicName: []byte(*topicName),
			Message:   []byte(*message),
		},
	}

	// Return the command.
	return cmd, nil
}
