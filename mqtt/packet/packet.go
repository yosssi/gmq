package packet

import "io"

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
	// Type return the MQTT Control Packet type of the Packet.
	Type() (byte, error)
}
