package mqtt

// Values of QoS.
const (
	// QoS0 represents "QoS 0: At most once delivery".
	QoS0 = iota
	// QoS1 represents "QoS 1: At least once delivery".
	QoS1
	// QoS2 represents "QoS 2: Exactly once delivery".
	QoS2
)
