package packet

// Base holds the common fields and method among MQTT Control Packets.
type Base struct {
	FixedHeader    []byte
	VariableHeader []byte
	Payload        []byte
}
