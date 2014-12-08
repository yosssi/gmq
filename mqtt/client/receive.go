package client

import (
	"bufio"
	"io"

	"github.com/yosssi/gmq/mqtt/packet"
)

// Receive receives an MQTT Control Packet from the Server.
func Receive(r *bufio.Reader) (byte, packet.Packet, error) {
	// Get the first byte of the Packet.
	b, err := r.ReadByte()
	if err != nil {
		return 0x00, nil, err
	}

	// Extract the MQTT Control Packet Type from the first byte.
	ptype := b >> 4

	// Create the Fixed header.
	fixedHeader := []byte{b}

	// Get and decode the Remaining Length.
	var mp uint32 = 1 // multiplier
	var rl uint32     // the Remaining Length
	for {
		b, err = r.ReadByte()
		if err != nil {
			return 0x00, nil, err
		}

		fixedHeader = append(fixedHeader, b)

		rl += uint32(b&0x7F) * mp

		if b&0x80 == 0 {
			break
		}

		mp *= 128
	}

	// Create the Remaining (the Variable header and the Payload).
	remaining := make([]byte, rl)

	if rl > 0 {
		if _, err = io.ReadFull(r, remaining); err != nil {
			return 0x00, nil, err
		}
	}

	var p packet.Packet

	switch ptype {
	case packet.TypeCONNACK:
		// Create the CONNACK Packet from the byte data to validate the data.
		if p, err = packet.NewCONNACKFromBytes(fixedHeader, remaining); err != nil {
			return 0x00, nil, err
		}
	}

	return ptype, p, nil
}
