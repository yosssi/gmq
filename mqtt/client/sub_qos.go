package client

// subQoS represents the QoS of the Subscription.
type subQoS struct {
	// reqQoS represents the QoS
	// which is requested by the Client.
	reqQoS byte
	// grantQoS represents the QoS
	// which is granted QoS by the Server.
	grantQoS byte
}
