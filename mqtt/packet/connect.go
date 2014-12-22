package packet

// connect represents a CONNECT Packet.
type connect struct {
	base
	// clientID is the Client Identifier of the payload.
	clientID []byte
	// userName is the User Name of the payload.
	userName []byte
	// password is the Password of the payload.
	password []byte
	// cleanSession is the Clean Session of the variable header.
	cleanSession bool
	// keepAlive is the Keep Alive of the variable header.
	keepAlive uint16
	// willTopic is the Will Topic of the payload.
	willTopic []byte
	// willMessage is the Will Message of the payload.
	willMessage []byte
	// willQoS is the Will QoS of the variable header.
	willQoS byte
	// willRetain is the Will Retain of the variable header.
	willRetain bool
}

// NewCONNECT creates and returns a CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}

	return nil, nil
}
