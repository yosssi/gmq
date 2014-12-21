package packet

import (
	"bytes"
	"io"
)

// Base holds the common fields and methods among MQTT Control Packets.
type Base struct {
	FixedHeader    []byte
	VariableHeader []byte
	Payload        []byte
}

// WriteTo writes the Packet data to the writer.
func (p *Base) WriteTo(w io.Writer) (int64, error) {
	var bf bytes.Buffer

	// Write the Fixed header, the Variable header and the Payload to the buffer.
	bf.Write(p.FixedHeader)
	bf.Write(p.VariableHeader)
	bf.Write(p.Payload)

	// Write the buffered data to the writer.
	n, err := w.Write(bf.Bytes())

	return int64(n), err
}

// Type returns the MQTT Control Packet type.
func (p *Base) Type() (byte, error) {
	return typeFromBytes(p.FixedHeader)
}

// appendRemainingLength appends the Remaining Length to the Fixed Header.
func (p *Base) appendRemainingLength() {
	// Calculate the Remaining Length.
	rl := encodeLength(uint(len(p.VariableHeader) + len(p.Payload)))

	// Append the Remaining Length to the slice and set it to the Fixed Header.
	p.FixedHeader = appendRemainingLength(p.FixedHeader, rl)
}

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
