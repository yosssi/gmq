package packet

import "io"

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
}

// encodeUint16 converts an unsigned 16-bit integer into
// a slice of bytes in big-endian order.
func encodeUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n)}
}

// encodeLength takes a uint value and returns uint32 of that value
// in the encoding mechanism defined for the MQTT remaining lengths.
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
