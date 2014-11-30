package client

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
func (opts *SessionOptions) Init() {
	// Initialize the CleanSession field
	if opts.CleanSession == nil {
		opts.CleanSession = &DefaultCleanSession
	}

	// Initialize the ClientID filed
	if opts.ClientID == "" {
		opts.ClientID = hostname
	}
}
