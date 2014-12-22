package packet

// CONNECTOptions represents options for a CONNECT Packet.
type CONNECTOptions struct {
	// ClientID is the Client Identifier of the payload.
	ClientID []byte
	// UserName is the User Name of the payload.
	UserName []byte
	// Password is the Password of the payload.
	Password []byte
	// CleanSession is the Clean Session of the variable header.
	CleanSession bool
	// KeepAlive is the Keep Alive of the variable header.
	KeepAlive uint
	// WillTopic is the Will Topic of the payload.
	WillTopic []byte
	// WillMessage is the Will Message of the payload.
	WillMessage []byte
	// WillQoS is the Will QoS of the variable header.
	WillQoS uint
	// WillRetain is the Will Retain of the variable header.
	WillRetain bool
}
