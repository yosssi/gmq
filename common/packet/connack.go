package packet

import "errors"

// Length of the Variable header of the CONNACK Packet.
const lenCONNACKVariableHeader = 2

// Error values
var ErrInvalidCONNACKVariableHeaderLen = errors.New("the length of the Variable header of the CONNACK Packet is invalid")

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
	// Check the length of the Variable header of the CONNACK Packet.
	if len(variableHeader) != lenCONNACKVariableHeader {
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
