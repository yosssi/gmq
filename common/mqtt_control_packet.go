package common

// MQTTControlPacket representa an MQTT Control Packet.
type MQTTControlPacket struct {
	// FixedHeader is a Fix header of the MQTT Control Packet.
	FixedHeader []byte
	// VariableHeader is a Variable header of the MQTT Control Packet.
	VariableHeader []byte
	// Payload is a Payload of the MQTT Control Packet.
	Payload []byte
}
