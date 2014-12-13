package main

import (
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
	sendc chan packet.Packet
	recvc chan packet.Packet
	errc  chan error
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli:   client.New(nil),
		sendc: make(chan packet.Packet, defaultSendcBufSize),
		recvc: make(chan packet.Packet, defaultRecvcBufSize),
		errc:  make(chan error, defaultErrcBufSize),
	}
}
