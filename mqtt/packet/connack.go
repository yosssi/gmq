package packet

import "errors"

// Length of the Fixed header of the CONNACK Packet.
const lenCONNACKFixedHeader = 2

// Length of the Variable header of the CONNACK Packet.
const lenCONNACKVariableHeader = 2

// First byte of the Fixed header of the CONNACK Packet.
const firstByteCONNACKFixedHeader = 0x20

// Error values
var (
	ErrInvalidCONNACKFixedHeaderLen       = errors.New("the length of the Fixed header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKVariableHeaderLen    = errors.New("the length of the Variable header of the CONNACK Packet is invalid")
	ErrInvalidfirstByteCONNACKFixedHeader = errors.New("the first byte of the Fixed header of the CONNACK Packet is invalid")
)

// CONNACK represents the CONNACK Packet.
type CONNACK struct {
	Base
	// SessionPresentFlag is the Session Present Flag of
	// the Connect Acknowledge Flags.
	SessionPresentFlag bool
	// ConnectReturnCode is the Connect Return Code.
	ConnectReturnCode byte
}

// NewCONNACKFromBytes creates the CONNACK Packet from the byte data and returns it.
func NewCONNACKFromBytes(fixedHeader, variableHeader []byte) (*CONNACK, error) {
	// Check the length of the Fixed header of the CONNACK Packet.
	if len(fixedHeader) != lenCONNACKFixedHeader {
		return nil, ErrInvalidCONNACKFixedHeaderLen
	}

	// Check the length of the Variable header of the CONNACK Packet.
	if len(variableHeader) != lenCONNACKVariableHeader {
		return nil, ErrInvalidCONNACKVariableHeaderLen
	}

	// Check the first byte of the Fixed header of the CONNACK Packet.
	if fixedHeader[0] != firstByteCONNACKFixedHeader {
		return nil, ErrInvalidCONNACKVariableHeaderLen
	}

	// Create a CONNACK Packet.
	p := &CONNACK{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = fixedHeader

	// Set the Variable header to the Packet.
	p.VariableHeader = variableHeader

	// Set the Session Present Flag to the Packet.
	p.SessionPresentFlag = (variableHeader[0] & 1) == 1

	// Set the Connect Return Code to the Packet.
	p.ConnectReturnCode = variableHeader[1]

	return p, nil
}
