package main

// commandDisconn represents a disconn command.
type commandDisconn struct {
	ctx *context
}

// run sends a DISCONNECT Packet to the Server and
// disconnects the Network Connection.
func (cmd *commandDisconn) run() error {
	return disconnect(cmd.ctx)
}

// newCommandDisconn creates and returns a disconn command.
func newCommandDisconn(ctx *context) *commandDisconn {
	return &commandDisconn{
		ctx: ctx,
	}
}

// disconnect disconnects the Network Connection.
func disconnect(ctx *context) error {
	// Close the Network Connection.
	if err := closeConn(ctx); err != nil {
		return err
	}

	// Wait until all goroutines which are accessing to the Network Connection end.
	ctx.wg.Wait()

	// Clear the Network Connection.
	clearConn(ctx)

	return nil
}

// closeConn closes the Network Connection.
func closeConn(ctx *context) error {
	// Lock for disconnecting the Network Connection.
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// Backup the current connected state.
	backedConnected := ctx.connected

	// Set the connected state false.
	ctx.connected = false

	// Disconnect the Network Connection.
	if err := ctx.cli.Disconnect(); err != nil {
		// Restore the connected state.
		ctx.connected = backedConnected

		// Return the error.
		return err
	}

	return nil
}

// clearConn clears the Netwrok Connection.
func clearConn(ctx *context) {
	// Lock for clearing the Network Connection.
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.cli.ClearConnection()
}
