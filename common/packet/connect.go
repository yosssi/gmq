package packet

import "io"

// connect represents a CONNECT Packet.
type connect struct {
	// CleanSession is the Clean Session of the connect flags.
	cleanSession bool
	// WillTopic is the Will Topic of the payload.
	willTopic string
	// WillMessage is the Will Message of the payload.
	willMessage string
	// WillQoS is the Will QoS of the connect flags.
	willQoS uint
	// WillRetain is the Will Retain of the connect flags.
	willRetain bool
	// UserName is the user name used by the server for authentication and authorization.
	userName string
	// Password is the password used by the server for authentication and authorization.
	password string
	// KeepAlive is the Keep Alive in the variable header.
	keepAlive uint
}

// WriteTo writes the Packet data to the writer.
func (p *connect) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

// NewCONNECT creates and returns a CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) Packet {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}
	opts.Init()

	return &connect{
		cleanSession: *opts.CleanSession,
		willTopic:    opts.WillTopic,
		willMessage:  opts.WillMessage,
		willQoS:      opts.WillQoS,
		willRetain:   opts.WillRetain,
		userName:     opts.UserName,
		password:     opts.Password,
		keepAlive:    *opts.KeepAlive,
	}
}
