package packet

import "errors"

// Length of the Fixed header of the CONNACK Packet
const lenCONNACKFixedHeader = 2

// Length of the Variable header of the CONNACK Packet
const lenCONNACKVariableHeader = 2

// First byte of the Fixed header of the CONNACK Packet
const firstByteCONNACKFixedHeader = TypeCONNACK << 4

// Connect Return code values of the CONNACK Packet
const (
	CONNACKConnectionAccepted                           = 0x00
	CONNACKConnectionRefusedUnacceptableProtocolVersion = 0x01
	CONNACKConnectionRefusedIdentifierRejected          = 0x02
	CONNACKConnectionRefusedServerUnavailable           = 0x03
	CONNACKConnectionRefusedBadUserNameOrPassword       = 0x04
	CONNACKConnectionRefusedNotAuthorized               = 0x05
)

// Error values
var (
	ErrInvalidCONNACKFixedHeaderLen          = errors.New("the length of the Fixed header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKVariableHeaderLen       = errors.New("the length of the Variable header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKFixedHeaderFirstByte    = errors.New("the first byte of the Fixed header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKRemainingLength         = errors.New("the Remaining Length of the Fixed header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKConnectAcknowledgeFlags = errors.New("the Connect Acknowledge Flags of the Variable header of the CONNACK Packet is invalid")
	ErrInvalidCONNACKConnectReturnCode       = errors.New("the Connect Return code of the Variable header of the CONNACK Packet is invalid")
)

// Valid Connect Return code values
var validConnectReturnCodes = []byte{
	CONNACKConnectionAccepted,
	CONNACKConnectionRefusedUnacceptableProtocolVersion,
	CONNACKConnectionRefusedIdentifierRejected,
	CONNACKConnectionRefusedServerUnavailable,
	CONNACKConnectionRefusedBadUserNameOrPassword,
	CONNACKConnectionRefusedNotAuthorized,
}

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
	// Check the length of the Fixed header.
	if len(fixedHeader) != lenCONNACKFixedHeader {
		return nil, ErrInvalidCONNACKFixedHeaderLen
	}

	// Check the length of the Variable header.
	if len(variableHeader) != lenCONNACKVariableHeader {
		return nil, ErrInvalidCONNACKVariableHeaderLen
	}

	// Check the first byte of the Fixed header.
	if fixedHeader[0] != firstByteCONNACKFixedHeader {
		return nil, ErrInvalidCONNACKFixedHeaderFirstByte
	}

	// Check the Remaining Length of the Fixed header.
	if fixedHeader[1] != lenCONNACKVariableHeader {
		return nil, ErrInvalidCONNACKRemainingLength
	}

	// Check the Connect Acknowledge Flags of the Variable header.
	if variableHeader[0]>>1 != 0x00 {
		return nil, ErrInvalidCONNACKConnectAcknowledgeFlags
	}

	// Check the Connect Return code of the Variable header.
	connectReturnCode := variableHeader[1]

	var validConnectReturnCode bool

	for _, c := range validConnectReturnCodes {
		if connectReturnCode == c {
			validConnectReturnCode = true
			break
		}
	}

	if !validConnectReturnCode {
		return nil, ErrInvalidCONNACKConnectReturnCode
	}

	// Create a CONNACK Packet.
	p := &CONNACK{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = fixedHeader

	// Set the Variable header to the Packet.
	p.VariableHeader = variableHeader

	// Set the Session Present Flag to the Packet.
	p.SessionPresentFlag = (variableHeader[0]<<7 == 0x80)

	// Set the Connect Return Code to the Packet.
	p.ConnectReturnCode = connectReturnCode

	// Return the CONNACK Packet.
	return p, nil
}
