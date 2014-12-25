package packet

// puback represents a PUBACK Packet.
type pubrel struct {
	base
	// packetID is the Packet Identifier of the variable header.
	packetID uint16
}

// PacketID returns the Packet Identifier of the Packet.
func (p *pubrel) PacketID() uint16 {
	return p.packetID
}

// setFixedHeader sets the fixed header to the Packet.
func (p *pubrel) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypePUBREL<<4|0x02)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *pubrel) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.packetID)...)
}

// NewPUBREL creates and returns a PUBREL Packet.
func NewPUBREL(opts *PUBRELOptions) Packet {
	// Initialize the options.
	if opts == nil {
		opts = &PUBRELOptions{}
	}

	// Create a PUBREL Packet.
	p := &pubrel{
		packetID: opts.PacketID,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p
}
