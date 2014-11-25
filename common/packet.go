package common

// Packet represent an MQTT Control Packet.
type Packet struct {
	// FixedHeader is a Fix header of the MQTT Control Packet.
	FixedHeader []byte
	// VariableHeader is a Variable header of the MQTT Control Packet.
	VariableHeader []byte
	// Payload is a Payload of the MQTT Control Packet.
	Payload []byte
}

// NewPacketCONNECT creates and returns a CONNECT Packet.
func NewPacketCONNECT(opts *OptionsPacketCONNECT) *Packet {
	// Initialize the options.
	if opts == nil {
		opts = &OptionsPacketCONNECT{}
	}
	opts.Init()

	return nil
}
