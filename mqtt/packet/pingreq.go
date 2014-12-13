package packet

// PINGREQ represents the PINGREQ Packet.
type PINGREQ struct {
	Base
}

// NewPINGREQ creates and returns the PINGREQ Packet.
func NewPINGREQ() Packet {
	// Create the PINGREQ Packet.
	p := &PINGREQ{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = []byte{TypePINGREQ << 4, 0x00}

	// Return the Packet.
	return p
}
