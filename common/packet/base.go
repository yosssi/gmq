package packet

// base holds the common fields and method among MQTT Control Packets.
type base struct {
	fixedHeader    []byte
	variableHeader []byte
	payload        []byte
}
