package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
)

// context represents a context of GMQ Client.
type context struct {
	cli     *client.Client
	closing bool
	mu      sync.RWMutex
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli: client.New(nil),
	}
}
