package packet

// appendRemainingLength append the Remaining Length to the slice
// and returns it.
func appendRemainingLength(b []byte, rl uint32) []byte {
	switch {
	case rl&0xFF000000 > 0:
		b = append(b, byte((rl&0xFF000000)>>24))
		b = append(b, byte((rl&0x00FF0000)>>16))
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	case rl&0x00FF0000 > 0:
		b = append(b, byte((rl&0x00FF0000)>>16))
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	case rl&0x0000FF00 > 0:
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	default:
		b = append(b, byte(rl&0x000000FF))
	}

	return b
}
