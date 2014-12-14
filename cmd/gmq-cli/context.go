package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Default buffer size of the channels
const (
	defaultSendcBufSize = 1024
	defaultRecvcBufSize = 1024
	defaultErrcBufSize  = 1024
)

// context represents a context of GMO Client.
type context struct {
	cli   *client.Client
	climu *sync.RWMutex
	sendc chan packet.Packet
	recvc chan packet.Packet
	errc  chan error
}

// disconnect disconnects the Network Connection.
func (ctx *context) disconnect() error {
	ctx.climu.Lock()
	defer ctx.climu.Unlock()

	if err := ctx.cli.Disconnect(); err != nil {
		return err
	}

	return nil
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli:   client.New(nil),
		climu: new(sync.RWMutex),
		sendc: make(chan packet.Packet, defaultSendcBufSize),
		recvc: make(chan packet.Packet, defaultRecvcBufSize),
		errc:  make(chan error, defaultErrcBufSize),
	}
}
