package packet

import "errors"

// Length of the fixed header of the CONNACK Packet
const lenCONNACKFixedHeader = 2

// Length of the variable header of the CONNACK Packet
const lenCONNACKVariableHeader = 2

// Connect Return code values
const (
	Accepted                    byte = 0x00
	UnacceptableProtocolVersion byte = 0x01
	IdentifierRejected          byte = 0x02
	ServerUnavailable           byte = 0x03
	BadUserNameOrPassword       byte = 0x04
	NotAuthorized               byte = 0x05
)

// Error values
var (
	ErrInvalidFixedHeader       = errors.New("invalid fixed header")
	ErrInvalidVariableHeaderLen = errors.New("invalid length of the variable header")
	ErrInvalidVariableHeader    = errors.New("invalid variable header")
	ErrInvalidRemainingLength   = errors.New("invalid Remaining Length")
	ErrInvalidConnectReturnCode = errors.New("invalid Connect Return code")
)

// connack represents a CONNACK Packet.
type connack struct {
	base
	// sessionPresent is the Session Present of the variable header.
	sessionPresent bool
	// connectReturnCode is the Connect Return code of the variable header.
	connectReturnCode byte
}

// NewCONNACKFromBytes creates the CONNACK Packet
// from the byte data and returns it.
func NewCONNACKFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validateCONNACKBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Create a CONNACK Packet.
	p := &connack{
		sessionPresent:    variableHeader[0]<<7 == 0x80,
		connectReturnCode: variableHeader[1],
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validateCONNACKBytes validates the fixed header and the variable header.
func validateCONNACKBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Check the length of the fixed header.
	if len(fixedHeader) != lenCONNACKFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	if ptype != TypeCONNACK {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenCONNACKVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenCONNACKVariableHeader {
		return ErrInvalidVariableHeaderLen
	}

	// Check the reserved bits of the variable header.
	if variableHeader[0]>>1 != 0x00 {
		return ErrInvalidVariableHeader
	}

	// Check the Connect Return code of the variable header.
	switch variableHeader[1] {
	case
		Accepted,
		UnacceptableProtocolVersion,
		IdentifierRejected,
		ServerUnavailable,
		BadUserNameOrPassword,
		NotAuthorized:
	default:
		return ErrInvalidConnectReturnCode
	}

	return nil
}
