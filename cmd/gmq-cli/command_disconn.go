package main

// commandDisconn represents a disconn command.
type commandDisconn struct {
	ctx *context
}

// run sends a DISCONNECT Packet to the Server and
// disconnects the Network Connection.
func (cmd *commandDisconn) run() error {
	return disconnectWithLock(cmd.ctx)
}

// newCommandDisconn creates and returns a disconn command.
func newCommandDisconn(ctx *context) *commandDisconn {
	return &commandDisconn{
		ctx: ctx,
	}
}
