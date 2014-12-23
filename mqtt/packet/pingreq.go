package packet

// pingreq represents a PINGREQ Packet.
type pingreq struct {
	base
}

// NewPINGREQ creates and returns a PINGREQ Packet.
func NewPINGREQ() Packet {
	// Create a PINGREQ Packet.
	p := &pingreq{}

	// Set the fixed header to the Packet.
	p.fixedHeader = []byte{TypePINGREQ << 4, 0x00}

	// Return the Packet.
	return p
}
