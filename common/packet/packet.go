package packet

import (
	"encoding/binary"
	"io"
)

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
}

//encodeUint16 takes a uint16 and returns a slice of bytes
//representing its value in network order
func encodeUint16(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

//encodeLength takes a uint value and returns uint32 of that value
//in the encoding mechanism defined for the MQTT remaining lengths
func encodeLength(length uint) uint32 {
	value := uint32(0)
	digit := uint32(0)
	x := uint32(length)
	for x > 0 {
		if value != 0 {
			value <<= 8
		}
		digit = x % 128
		x /= 128
		if x > 0 {
			digit |= 0x80
		}
		value |= uint32(digit)
	}
	return value
}
