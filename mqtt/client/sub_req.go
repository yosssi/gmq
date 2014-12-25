package client

// SubReq represents subscription request.
type SubReq struct {
	// TopicFilter is the Topic Filter of the Subscription.
	TopicFilter string
	// QoS is the requsting QoS.
	QoS byte
}
