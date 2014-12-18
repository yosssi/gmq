package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

const sendBufSize = 1024

// context represents a context of GMQ Client.
type context struct {
	mu  sync.RWMutex
	cli *client.Client

	disconn chan struct{}

	wg         sync.WaitGroup
	connack    chan struct{}
	connackEnd chan struct{}
	send       chan packet.Packet
	sendEnd    chan struct{}
}

// initChan initializes the channels of the context.
func (ctx *context) initChan() {
	ctx.connack = make(chan struct{}, 1)
	ctx.connackEnd = make(chan struct{}, 1)
	ctx.send = make(chan packet.Packet, sendBufSize)
	ctx.sendEnd = make(chan struct{}, 1)
}

// newContext creates and returns a context.
func newContext() *context {
	ctx := &context{
		cli:     client.New(nil),
		disconn: make(chan struct{}, 1),
	}

	ctx.initChan()

	return ctx
}
