package packet

// encodeUint16 converts an unsigned 16-bit integer
// into a slice of bytes in big-endian order.
func encodeUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n)}
}

// encodeLength encodes an unsigned integer
// using a variable length encoding scheme.
func encodeLength(length uint) uint32 {
	var value, digit uint32

	n := uint32(length)

	for n > 0 {
		if value != 0 {
			value <<= 8
		}

		digit = n % 128

		n /= 128

		if n > 0 {
			digit |= 0x80
		}

		value |= uint32(digit)
	}

	return value
}
