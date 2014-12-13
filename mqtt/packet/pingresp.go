package packet

import "errors"

// Length of the Fixed header of the PINGRESP Packet
const lenPINGRESPFixedHeader = 2

// First byte of the Fixed header of the PINGRESP Packet
const firstBytePINGRESPFixedHeader = TypePINGRESP << 4

// Error values
var (
	ErrPINGRESPInvalidFixedHeaderLen       = errors.New("the length of the Fixed header of the PINGRESP Packet is invalid")
	ErrPINGRESPInvalidFixedHeaderFirstByte = errors.New("the first byte of the Fixed header of the PINGRESP Packet is invalid")
	ErrPINGRESPInvalidRemainingLength      = errors.New("the Remaining Length of the Fixed header of the PINGRESP Packet is invalid")
)

// PINGRESP represents the PINGRESP Packet.
type PINGRESP struct {
	Base
}

// NewPINGRESPFromBytes creates the PINGRESP Packet
// from the byte data and returns it.
func NewPINGRESPFromBytes(fixedHeader []byte) (Packet, error) {
	// Check the length of the Fixed header.
	if len(fixedHeader) != lenPINGRESPFixedHeader {
		return nil, ErrPINGRESPInvalidFixedHeaderLen
	}

	// Check the first byte of the Fixed header.
	if fixedHeader[0] != firstBytePINGRESPFixedHeader {
		return nil, ErrPINGRESPInvalidFixedHeaderFirstByte
	}

	// Check the Remaining Length of the Fixed header.
	if fixedHeader[1] != 0x00 {
		return nil, ErrPINGRESPInvalidRemainingLength
	}

	// Create the PINGRESP Packet.
	p := &PINGRESP{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = fixedHeader

	// Return the Packet.
	return p, nil
}
