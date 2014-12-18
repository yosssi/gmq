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
	connack chan struct{}
	send    chan packet.Packet
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli:     client.New(nil),
		disconn: make(chan struct{}, 1),
		connack: make(chan struct{}, 1),
		send:    make(chan packet.Packet, sendBufSize),
	}
}
