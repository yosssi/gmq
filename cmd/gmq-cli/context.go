package main

import (
	"errors"
	"sync"

	"github.com/yosssi/gmq/mqtt/client"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Buffer size of the send channel
const sendBufSize = 1024

// Max Packet Identifier
const maxPacketID = 65535

// Error value
var errPacketIDExhaused = errors.New("Packet Identifiers are exhausted")

// context represents a context of GMQ Client.
type context struct {
	mu            sync.RWMutex
	cli           *client.Client
	disconnecting bool

	wgMain     sync.WaitGroup
	disconn    chan struct{}
	disconnEnd chan struct{}

	wg         sync.WaitGroup
	connack    chan struct{}
	connackEnd chan struct{}
	send       chan packet.Packet
	sendEnd    chan struct{}

	// packetIDs holds Packet Identifiers currently used.
	muPacketIds sync.Mutex
	packetIDs   map[uint16]struct{}
}

// initChan initializes the channels of the context.
func (ctx *context) initChan() {
	ctx.connack = make(chan struct{}, 1)
	ctx.connackEnd = make(chan struct{}, 1)
	ctx.send = make(chan packet.Packet, sendBufSize)
	ctx.sendEnd = make(chan struct{}, 1)
}

// initPacketIDs initializes packetIDs.
func (ctx *context) initPacketIDs() {
	// Lock for the initialization.
	ctx.muPacketIds.Lock()
	defer ctx.muPacketIds.Unlock()

	ctx.packetIDs = make(map[uint16]struct{})
}

// generatePacketID generates and returns a Packet Identifier.
func (ctx *context) generatePacketID() (uint16, error) {
	// Lock for the generation of the Packet Identifier.
	ctx.muPacketIds.Lock()
	defer ctx.muPacketIds.Unlock()

	var id uint16
	for {
		// Find an id which does not exist in packetIDs.
		if _, exist := ctx.packetIDs[id]; !exist {
			// Set the id to packetIDs.
			ctx.packetIDs[id] = struct{}{}
			// Return the id.
			return id, nil
		}

		if id == maxPacketID {
			break
		}

		id++
	}

	// Return an error if available ids are not found in packetIDs.
	return 0, errPacketIDExhaused
}

// newContext creates and returns a context.
func newContext() *context {
	ctx := &context{
		cli:        client.New(nil),
		disconn:    make(chan struct{}, 1),
		disconnEnd: make(chan struct{}, 1),
	}

	ctx.initChan()

	ctx.initPacketIDs()

	return ctx
}
