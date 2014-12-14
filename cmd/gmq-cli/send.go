package main

import "github.com/yosssi/gmq/mqtt/packet"

// sendWithLock locks the Client and sends a Packet to the Server.
func sendWithLock(ctx *context, p packet.Packet) error {
	ctx.climu.Lock()
	defer ctx.climu.Unlock()

	return ctx.cli.Send(p)
}
