package packet

import "github.com/yosssi/gmq/mqtt"

// PUBLISH represents a PUBLISH Packet.
type PUBLISH struct {
	base
	// dup is the DUP flag of the fixed header.
	dup bool
	// qos is the QoS of the fixed header.
	qos byte
	// retain is the Retain of the fixed header.
	retain bool
	// topicName is the Topic Name of the varible header.
	topicName []byte
	// packetID is the Packet Identifier of the variable header.
	PacketID uint16
	// message is the Application Message of the payload.
	message []byte
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBLISH) setFixedHeader() {
	// Define the first byte of the fixed header.
	b := TypePUBLISH << 4

	// Set 1 to the Bit 3 if the DUP flag is true.
	if p.dup {
		b |= 0x08
	}

	// Set the value of the Will QoS to the Bit 2 and 1.
	b |= p.qos << 1

	// Set 1 to the Bit 0 if the Retain is true.
	if p.retain {
		b |= 0x01
	}

	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, b)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBLISH) setVariableHeader() {
	// Append the Topic Name to the variable header.
	p.variableHeader = appendLenStr(p.variableHeader, p.topicName)

	if p.qos != mqtt.QoS0 {
		// Append the Packet Identifier to the variable header.
		p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
	}
}

// setPayload sets the payload to the Packet.
func (p *PUBLISH) setPayload() {
	p.payload = p.message
}

// NewPUBLISH creates and returns a PUBLISH Packet.
func NewPUBLISH(opts *PUBLISHOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &PUBLISHOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a PUBLISH Packet.
	p := &PUBLISH{
		dup:       opts.DUP,
		qos:       opts.QoS,
		retain:    opts.Retain,
		topicName: opts.TopicName,
		PacketID:  opts.PacketID,
		message:   opts.Message,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the payload to the Packet.
	p.setPayload()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}
