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
func (p Base) WriteTo(w io.Writer) (int64, error) {
	var bf bytes.Buffer

	// Write the Fixed header, the Variable header and the Payload to the buffer.
	bf.Write(p.FixedHeader)
	bf.Write(p.VariableHeader)
	bf.Write(p.Payload)

	// Write the buffered data to the writer.
	n, err := w.Write(bf.Bytes())

	return int64(n), err
}
