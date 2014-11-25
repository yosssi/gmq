package packet

import "io"

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
}
