package mqtt

import "os"

// Defalut value
var DefaultCleanSession = true

// Hostname
var hostname, _ = os.Hostname()

// SessionOptions represents options for a Session.
type SessionOptions struct {
	CleanSession *bool
	ClientID     string
}

// Init initializes the options.
func (o *SessionOptions) Init() {
	// Initialize the CleanSession field
	if o.CleanSession == nil {
		o.CleanSession = &DefaultCleanSession
	}

	// Initialize the ClientID filed
	if o.ClientID == "" {
		o.ClientID = hostname
	}
}
