package packet

import (
	"errors"
	"io"
)

// Error value
var ErrInvalidPacketType = errors.New("invalid MQTT Control Packet type")

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
	// Type returns the MQTT Control Packet type of the Packet.
	Type() (byte, error)
	// PacketID returns the Packet Identifier of the Packet.
	PacketID() uint16
}

// NewFromBytes creates a Packet from the byte data and returns it.
func NewFromBytes(fixedHeader FixedHeader, remaining []byte) (Packet, error) {
	// Extract the MQTT Control Packet type from the fixed header.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return nil, err
	}

	// Create and return a Packet.
	switch ptype {
	case TypeCONNACK:
		return NewCONNACKFromBytes(fixedHeader, remaining)
	case TypePUBACK:
		return NewPUBACKFromBytes(fixedHeader, remaining)
	case TypePUBREC:
		return NewPUBRECFromBytes(fixedHeader, remaining)
	case TypePINGRESP:
		return NewPINGRESPFromBytes(fixedHeader, remaining)
	default:
		return nil, ErrInvalidPacketType
	}
}
