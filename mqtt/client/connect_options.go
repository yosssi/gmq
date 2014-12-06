package client

// Default values
const (
	DefaultNetwork = "tcp"
	DefaultAddress = "localhost:1883"
)

// ConnectOptions is options for the Client's Connect method.
type ConnectOptions struct {
	// Network is the network which the Client connects to.
	Network string
	// Address is the address which the Client connects to.
	Address string
}

// Init initializes the ConnectOptions.
func (opts *ConnectOptions) Init() {
	if opts.Network == "" {
		opts.Network = DefaultNetwork
	}

	if opts.Address == "" {
		opts.Address = DefaultAddress
	}
}
