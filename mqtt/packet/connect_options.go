package packet

import "os"

// Defalut values
var (
	DefaultCleanSession      = true
	DefaultKeepAlive    uint = 60
)

// Hostname
var hostname, _ = os.Hostname()

// CONNECTOptions represents options for a CONNECT Packet.
type CONNECTOptions struct {
	// ClientID is the Client Identifier which identifies the Client to the Server.
	ClientID string
	// CleanSession is the Clean Session of the Connect Flags.
	CleanSession *bool
	// WillTopic is the Will Topic of the Payload.
	WillTopic string
	// WillMessage is the Will Message of the Payload.
	WillMessage string
	// WillQoS is the Will QoS of the Connect Flags.
	WillQoS uint
	// WillRetain is the Will Retain of the Connect Flags.
	WillRetain bool
	// UserName is the User Name used by the Server for authentication and authorization.
	UserName string
	// Password is the Password used by the Server for authentication and authorization.
	Password string
	// KeepAlive is the Keep Alive in the Variable header.
	KeepAlive *uint
}

// Init initializes the CONNECTOptions.
func (opts *CONNECTOptions) Init() {
	if opts.ClientID == "" {
		opts.ClientID = hostname
	}

	if opts.CleanSession == nil {
		opts.CleanSession = &DefaultCleanSession
	}

	if opts.KeepAlive == nil {
		opts.KeepAlive = &DefaultKeepAlive
	}
}
