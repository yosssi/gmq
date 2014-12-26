package packet

// PUBREL represents a PUBREL Packet.
type PUBREL struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBREL) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypePUBREL<<4|0x02)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBREL) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// NewPUBREL creates and returns a PUBREL Packet.
func NewPUBREL(opts *PUBRELOptions) Packet {
	// Initialize the options.
	if opts == nil {
		opts = &PUBRELOptions{}
	}

	// Create a PUBREL Packet.
	p := &PUBREL{
		PacketID: opts.PacketID,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p
}
