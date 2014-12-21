package packet

import (
	"errors"
	"strings"

	"github.com/yosssi/gmq/mqtt"
)

// Error values
var (
	ErrPUBLISHInvalidQoS                = errors.New("the QoS is invalid")
	ErrPUBLISHTopicNameContainsWildcard = errors.New("the Topic Name must not contain wildcard characters")
)

// PUBLISH represents the PUBLISH Packet.
type PUBLISH struct {
	Base
	// DUP is the DUP flag of the Fixed header.
	DUP bool
	// QoS is the QoS of the Fixed header.
	QoS uint
	// Retain is the Retain of the Fixed header.
	Retain bool
	// TopicName is the Topic Name of the Variable header.
	TopicName string
	// PacketID is the Packet Identifier of the Variable header.
	PacketID uint16
	// Message is the Application Message of the Payload.
	Message string
}

// setFixedHeader sets the Fixed header to the Packet.
func (p *PUBLISH) setFixedHeader() {
	// Create the first bit.
	var b byte

	b = TypePUBLISH << 4

	if p.DUP {
		b |= 0x80
	}

	b |= byte(p.QoS) << 1

	if p.Retain {
		b |= 0x01
	}

	// Append the first bit to the Fixed header.
	p.FixedHeader = append(p.FixedHeader, b)

	// Append the Remaining Length to the Fixed Header.
	p.appendRemainingLength()
}

// setVariableHeader sets the Variable header to the Packet.
func (p *PUBLISH) setVariableHeader() {
	topicName := []byte(p.TopicName)

	p.VariableHeader = append(p.VariableHeader, encodeUint16(uint16(len(topicName)))...)
	p.VariableHeader = append(p.VariableHeader, topicName...)
	p.VariableHeader = append(p.VariableHeader, encodeUint16(p.PacketID)...)
}

// setPayload sets the Payload to the Packet.
func (p *PUBLISH) setPayload() {
	p.Payload = []byte(p.Message)
}

// NewPUBLISH creates and returns a PUBLISH Packet.
func NewPUBLISH(opts *PUBLISHOptions) (*PUBLISH, error) {
	if opts == nil {
		opts = &PUBLISHOptions{}
	}

	// Check the QoS.
	if opts.QoS != mqtt.QoS0 && opts.QoS != mqtt.QoS1 && opts.QoS != mqtt.QoS2 {
		return nil, ErrPUBLISHInvalidQoS
	}

	// Check the Topic Name.
	if strings.Index(opts.TopicName, "*") != -1 {
		return nil, ErrPUBLISHTopicNameContainsWildcard
	}

	// Create a PUBLISH Packet.
	p := &PUBLISH{
		DUP:       opts.DUP,
		QoS:       opts.QoS,
		Retain:    opts.Retain,
		TopicName: opts.TopicName,
		PacketID:  opts.PacketID,
		Message:   opts.Message,
	}

	// Set the Variable header to the Packet.
	p.setVariableHeader()

	// Set the Payload to the Packet.
	p.setPayload()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}
