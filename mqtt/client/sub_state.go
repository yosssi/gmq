package client

// SubState represents the state of the Subscription.
type SubState struct {
	// ReqQoS represents the QoS
	// which is requested by the Client.
	ReqQoS byte
	// GrantQoS represents the QoS
	// which is granted QoS by the Server.
	GrantQoS byte
	// Acked represents if the SUBACK Packet was
	// sent from the Server or not.
	Acked bool
	// Failed represents if the Return Code represented
	// failure or not.
	Failed bool
}
