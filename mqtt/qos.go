package mqtt

// Values of QoS.
const (
	// QoS0 represents "QoS 0: At most once delivery".
	QoS0 uint = iota
	// QoS1 represents "QoS 1: At least once delivery".
	QoS1
	// QoS2 represents "QoS 2: Exactly once delivery".
	QoS2
)

// ValidQoS returns true if the input QoS is Qos0, QoS1 or QoS2.
func ValidQoS(qos uint) bool {
	return qos == QoS0 || qos == QoS1 || qos == QoS2
}
