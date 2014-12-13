package packet

import (
	"errors"
	"io"
)

// Error values
var (
	ErrInvalidFixedHeaderLen = errors.New("the length of the Fixed header is invalid")
	ErrInvalidPacketType     = errors.New("invalid MQTT Control Packet type")
)

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
	// Type returns the MQTT Control Packet type.
	Type() (byte, error)
}

// NewFromBytes creates a Packet from the byte data and returns it.
func NewFromBytes(fixedHeader, remaining []byte) (Packet, error) {
	ptype, err := typeFromBytes(fixedHeader)
	if err != nil {
		return nil, err
	}

	var p Packet

	switch ptype {
	case TypeCONNACK:
		// Create the CONNACK Packet from the byte data.
		if p, err = NewCONNACKFromBytes(fixedHeader, remaining); err != nil {
			return nil, err
		}
	case TypePINGRESP:
		// Create the PINGRESP Packet from the byte data.
		if p, err = NewPINGRESPFromBytes(fixedHeader); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidPacketType
	}

	// Ruturn the Packet.
	return p, nil
}

// typeFromBytes returns the MQTT Control Packet type.
func typeFromBytes(fixedHeader []byte) (byte, error) {
	if len(fixedHeader) < 1 {
		return 0x00, ErrInvalidFixedHeaderLen
	}

	return fixedHeader[0] >> 4, nil
}
