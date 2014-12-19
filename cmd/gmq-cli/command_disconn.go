package main

import (
	"errors"

	"github.com/yosssi/gmq/mqtt/client"
)

// Error value
var errDisconnSig = errors.New("disconnect signal could not be sent to the channel")

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
	// Lock for the disconnection.
	ctx.mu.Lock()

	// Disconnect the Network Connection.
	if err := ctx.cli.Disconnect(); err != nil {
		if err == client.ErrNotYetConnected {
			// Unlock.
			ctx.mu.Unlock()

			return err
		}

		// Close the Network Connection directly because
		// sending a DISCONNECT Packet to the Server failed.
		if err := ctx.cli.Close(); err != nil {
			printError(err, true)
		}
	}

	// Set the disconnecting flag true.
	ctx.disconnecting = true

	// Unlock.
	ctx.mu.Unlock()

	// Send the end signals to the channels.
	select {
	case ctx.connackEnd <- struct{}{}:
	default:
	}

	select {
	case ctx.sendEnd <- struct{}{}:
	default:
	}

	// Wait until all goroutines end.
	ctx.wg.Wait()

	// Lock for clearance of the Network Connection.
	ctx.mu.Lock()

	// Clear the Network Connection.
	ctx.cli.ClearConnection()

	// Initialize the channels of the context.
	ctx.initChan()

	// Set the disconnecting flag false.
	ctx.disconnecting = true

	// Unlock.
	ctx.mu.Unlock()

	return nil
}
