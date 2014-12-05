package packet

// DISCONNECT represents the DISCONNECT Packet.
type DISCONNECT struct {
	Base
}

// NewDISCONNECT creates and returns the DISCONNECT Packet.
func NewDISCONNECT() *DISCONNECT {
	// Create a DISCONNECT Packet.
	p := &DISCONNECT{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = []byte{TypeDISCONNECT << 4, 0x00}

	// Return the Packet.
	return p
}
