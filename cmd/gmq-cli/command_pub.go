package main

import (
	"flag"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

// commandPub represents a pub command.
type commandPub struct {
	ctx         *context
	publishOpts *packet.PUBLISHOptions
}

func (cmd *commandPub) run() error {
	// Lock for reading.
	cmd.ctx.mu.RLock()
	defer cmd.ctx.mu.RUnlock()

	// Check the existence of the Network Connection.
	if !cmd.ctx.cli.Connected() {
		return client.ErrNotYetConnected
	}

	// Create a PUBLISH Packet.
	p, err := packet.NewPUBLISH(cmd.publishOpts)
	if err != nil {
		return err
	}

	cmd.ctx.send <- p

	return nil
}

// newCommandPub creates and returns a pub command.
func newCommandPub(args []string, ctx *context) (*commandPub, error) {
	// Create a flag set.
	var flg flag.FlagSet

	// Define the flags.
	dup := flg.Bool("d", false, "DUP flag")
	qos := flg.Uint("q", mqtt.QoS0, "QoS")
	retain := flg.Bool("r", false, "Retain")
	topicName := flg.String("t", "", "Topic Name")
	message := flg.String("m", "", "Application Message")

	// Parse the flag.
	if err := flg.Parse(args); err != nil {
		return nil, errCmdArgsParse
	}

	// Generate a Packet Identifier.
	packetID, err := ctx.generatePacketID()
	if err != nil {
		return nil, err
	}

	// Create a conn command.
	cmd := &commandPub{
		ctx: ctx,
		publishOpts: &packet.PUBLISHOptions{
			DUP:       *dup,
			QoS:       *qos,
			Retain:    *retain,
			TopicName: *topicName,
			PacketID:  packetID,
			Message:   *message,
		},
	}

	return cmd, nil
}
