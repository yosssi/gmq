package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
)

// context represents a context of GMQ Client.
type context struct {
	mu        sync.RWMutex
	cli       *client.Client
	connected bool

	wg sync.WaitGroup
}

// isConnected returns the value of the connected state.
func (ctx *context) isConnected() bool {
	// Lock for reading.
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	// Return the value of the connected state.
	return ctx.connected
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli: client.New(nil),
	}
}
