package client

// Options represents options for an MQTT client.
type Options struct {
	// Host is the name of the host to connect to.
	Host string
	// Port is the port number of the host to connect to.
	Port uint16
	// CleanSession is a flag which corresponds to the CleanSession of the Connect Flags
	// in the CONNECT Packet.
	CleanSession *bool
}
