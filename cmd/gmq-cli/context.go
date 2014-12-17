package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
)

// Values of connState
const (
	connStateClosed uint32 = iota
	connStateOpen
	connStateDisconnecting
)

// context represents a context of GMQ Client.
type context struct {
	cli *client.Client

	connMu    sync.Mutex
	connState uint32
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli: client.New(nil),
	}
}
