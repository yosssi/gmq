package main

import (
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Default buffer size of the channels
const (
	defaultSendcBufSize    = 1024
	defaultRecvcBufSize    = 1024
	defaultErrcBufSize     = 1024
	defaultDisconncBufSize = 1024
)

// context represents a context of GMO Client.
type context struct {
	cli   *client.Client
	climu *sync.RWMutex

	sendc      chan packet.Packet
	sendEndc   chan struct{}
	sendEndedc chan struct{}

	readEndedc chan struct{}

	recvc      chan packet.Packet
	recvEndc   chan struct{}
	recvEndedc chan struct{}

	connackc      chan struct{}
	connackEndc   chan struct{}
	connackEndedc chan struct{}

	errc      chan error
	errEndc   chan struct{}
	errEndedc chan struct{}

	disconnc      chan struct{}
	disconnEndc   chan struct{}
	disconnEndedc chan struct{}
}

// newContext creates and returns a context.
func newContext() *context {
	return &context{
		cli:   client.New(nil),
		climu: new(sync.RWMutex),

		sendc:      make(chan packet.Packet, defaultSendcBufSize),
		sendEndc:   make(chan struct{}),
		sendEndedc: make(chan struct{}),

		readEndedc: make(chan struct{}),

		recvc:      make(chan packet.Packet, defaultRecvcBufSize),
		recvEndc:   make(chan struct{}),
		recvEndedc: make(chan struct{}),

		connackc:      make(chan struct{}),
		connackEndc:   make(chan struct{}),
		connackEndedc: make(chan struct{}),

		errc:      make(chan error, defaultErrcBufSize),
		errEndc:   make(chan struct{}),
		errEndedc: make(chan struct{}),

		disconnc:      make(chan struct{}, defaultDisconncBufSize),
		disconnEndc:   make(chan struct{}),
		disconnEndedc: make(chan struct{}),
	}
}
