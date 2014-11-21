package client

// Defalut values
const (
	DefaultHost = "localhost"
)

// Defalut value
var (
	DefaultPort         uint16 = 1883
	DefaultCleanSession        = true
)

// ConnOpts is options for the Client.Conn method.
type ConnOpts struct {
	// Host is the host name of the server to connect to.
	Host string
	// Port is the port number of the server to connect to.
	Port *uint16
	// CleanSession is the Clean Session of the connect flags.
	CleanSession *bool
	// WillTopic is the Will Topic of the payload.
	WillTopic string
	// WillMessage is the Will Message of the payload.
	WillMessage string
	// WillQoS is the Will QoS of the connect flags.
	WillQoS uint8
	// WillRetain is the Will Retain of the connect flags.
	WillRetain bool
	// UserName is the user name used by the server for authentication and authorization.
	UserName string
	// Password is the password used by the server for authentication and authorization.
	Password string
	// KeepAlive is the Keep Alive in the variable header.
	KeepAlive uint16
}

// Init initialize the ConnOpts.
func (opts *ConnOpts) Init() {
	if opts.Host == "" {
		opts.Host = DefaultHost
	}

	if opts.Port == nil {
		opts.Port = &DefaultPort
	}

	if opts.CleanSession == nil {
		opts.CleanSession = &DefaultCleanSession
	}
}
