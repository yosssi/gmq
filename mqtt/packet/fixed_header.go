package packet

import "errors"

// Error value
var ErrInvalidFixedHeaderLen = errors.New("the length of the fixed header is invalid")

// fixedHeader represents the fixed header of the Packet.
type fixedHeader []byte

// ptype extracts the MQTT Control Packet type from
// the fixed header and returns it.
func (fh fixedHeader) ptype() (byte, error) {
	// Check the length of the fixed header.
	if len(fh) < 1 {
		return 0x00, ErrInvalidFixedHeaderLen
	}

	// Extract the MQTT Control Packet type from
	// the fixed header and return it.
	return fh[0] >> 4, nil
}
