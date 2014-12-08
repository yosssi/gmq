package client

import "time"

// Default values
var DefaultConnTimeout = 60 * time.Second

// Options is options for the Client.
type Options struct {
	// ConnTimeout represents the timeout of the Network Connection.
	// 0 means that there is no timeout.
	ConnTimeout *time.Duration
}

// Init initializes the Options.
func (opts *Options) Init() {
	if opts.ConnTimeout == nil {
		opts.ConnTimeout = &DefaultConnTimeout
	}
}
