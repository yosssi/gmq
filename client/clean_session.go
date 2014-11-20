package client

// Values of a cleanSession
var (
	CleanSessionDisabled = newCleanSession(false)
	CleanSessionEnabled  = newCleanSession(true)
)

// cleanSession represents a flag corresponds to the CleanSession of the Connect Flags
// in the CONNECT Packet.
type cleanSession bool

// newCleanSession creates and returns a cleanSession.
func newCleanSession(b bool) *cleanSession {
	cs := cleanSession(b)
	return &cs
}
