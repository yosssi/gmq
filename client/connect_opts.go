package client

// Defalut values
const (
	DefaultHost = "localhost"
)

// Defalut value
var (
	DefaultPort         uint = 1883
	DefaultCleanSession      = true
	DefaultKeepAlive    uint = 60
)

// ConnectOpts is options for the Client.Conn method.
type ConnectOpts struct {
	// Host is the host name of the server to connect to.
	Host string
	// Port is the port number of the server to connect to.
	Port *uint
	// CleanSession is the Clean Session of the connect flags.
	CleanSession *bool
	// WillTopic is the Will Topic of the payload.
	WillTopic string
	// WillMessage is the Will Message of the payload.
	WillMessage string
	// WillQoS is the Will QoS of the connect flags.
	WillQoS uint
	// WillRetain is the Will Retain of the connect flags.
	WillRetain bool
	// UserName is the user name used by the server for authentication and authorization.
	UserName string
	// Password is the password used by the server for authentication and authorization.
	Password string
	// KeepAlive is the Keep Alive in the variable header.
	KeepAlive uint
}

// Init initialize the ConnectOpts.
func (opts *ConnectOpts) Init() {
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
