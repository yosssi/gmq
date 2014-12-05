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
	CONNACKAccepted                    = 0x00
	CONNACKUnacceptableProtocolVersion = 0x01
	CONNACKIdentifierRejected          = 0x02
	CONNACKServerUnavailable           = 0x03
	CONNACKBadUserNameOrPassword       = 0x04
	CONNACKNotAuthorized               = 0x05
)

// Error values
var (
	ErrCONNACKInvalidFixedHeaderLen          = errors.New("the length of the Fixed header of the CONNACK Packet is invalid")
	ErrCONNACKInvalidVariableHeaderLen       = errors.New("the length of the Variable header of the CONNACK Packet is invalid")
	ErrCONNACKInvalidFixedHeaderFirstByte    = errors.New("the first byte of the Fixed header of the CONNACK Packet is invalid")
	ErrCONNACKInvalidRemainingLength         = errors.New("the Remaining Length of the Fixed header of the CONNACK Packet is invalid")
	ErrCONNACKInvalidConnectAcknowledgeFlags = errors.New("the Connect Acknowledge Flags of the Variable header of the CONNACK Packet is invalid")
	ErrCONNACKInvalidConnectReturnCode       = errors.New("the Connect Return code of the Variable header of the CONNACK Packet is invalid")
)

// Valid Connect Return code values
var validConnectReturnCodes = []byte{
	CONNACKAccepted,
	CONNACKUnacceptableProtocolVersion,
	CONNACKIdentifierRejected,
	CONNACKServerUnavailable,
	CONNACKBadUserNameOrPassword,
	CONNACKNotAuthorized,
}

// CONNACK represents the CONNACK Packet.
type CONNACK struct {
	Base
	SessionPresent    bool
	ConnectReturnCode byte
}

// NewCONNACKFromBytes creates the CONNACK Packet from the byte data and returns it.
func NewCONNACKFromBytes(fixedHeader, variableHeader []byte) (*CONNACK, error) {
	// Check the length of the Fixed header.
	if len(fixedHeader) != lenCONNACKFixedHeader {
		return nil, ErrCONNACKInvalidFixedHeaderLen
	}

	// Check the length of the Variable header.
	if len(variableHeader) != lenCONNACKVariableHeader {
		return nil, ErrCONNACKInvalidVariableHeaderLen
	}

	// Check the first byte of the Fixed header.
	if fixedHeader[0] != firstByteCONNACKFixedHeader {
		return nil, ErrCONNACKInvalidFixedHeaderFirstByte
	}

	// Check the Remaining Length of the Fixed header.
	if fixedHeader[1] != lenCONNACKVariableHeader {
		return nil, ErrCONNACKInvalidRemainingLength
	}

	// Check the Connect Acknowledge Flags of the Variable header.
	if variableHeader[0]>>1 != 0x00 {
		return nil, ErrCONNACKInvalidConnectAcknowledgeFlags
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
		return nil, ErrCONNACKInvalidConnectReturnCode
	}

	// Create a CONNACK Packet.
	p := &CONNACK{}

	// Set the Fixed header to the Packet.
	p.FixedHeader = fixedHeader

	// Set the Variable header to the Packet.
	p.VariableHeader = variableHeader

	// Set the Session Present flag to the Packet.
	p.SessionPresent = (variableHeader[0]<<7 == 0x80)

	// Set the Connect Return Code to the Packet.
	p.ConnectReturnCode = connectReturnCode

	// Return the CONNACK Packet.
	return p, nil
}
