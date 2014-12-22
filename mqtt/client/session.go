package client

// session represents a Session which is a stateful interaction between a Client and a Server.
type session struct {
	// cleanSession is the Clean Session.
	cleanSession bool
	// clientID is the Client Identifier.
	clientID []byte
}
