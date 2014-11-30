package mqtt

// Session represents a Session which is a stateful interaction between a Client and a Server.
type Session struct {
	CleanSession bool
	ClientID     string
}

// NewSession creates and returns a Session.
func NewSession(opts *SessionOptions) *Session {
	// Initialize the options.
	if opts == nil {
		opts = &SessionOptions{}
	}
	opts.Init()

	// Create and return a Session.
	return &Session{
		CleanSession: *opts.CleanSession,
		ClientID:     opts.ClientID,
	}
}
