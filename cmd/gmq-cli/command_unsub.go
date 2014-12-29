package main

import (
	"flag"

	"github.com/yosssi/gmq/mqtt/client"
)

// Command name
const cmdNameUnsub = "unsub"

// commandUnsub represents an unsub command.
type commandUnsub struct {
	cli             *client.Client
	unsubscribeOpts *client.UnsubscribeOptions
}

// run sends an UNSUBSCRIBE Packet to the Server.
func (cmd *commandUnsub) run() error {
	return cmd.cli.Unsubscribe(cmd.unsubscribeOpts)
}

// newCommandUnsub creates and returns an unsub command.
func newCommandUnsub(args []string, cli *client.Client) (command, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	topicFilter := flg.String("t", "", "Topic Filter")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	// Create an unsub command.
	cmd := &commandUnsub{
		cli: cli,
		unsubscribeOpts: &client.UnsubscribeOptions{
			TopicFilters: [][]byte{
				[]byte(*topicFilter),
			},
		},
	}

	// Return the command.
	return cmd, nil
}
