package common

// Values of QoS.
const (
	QoS0 = iota
	QoS1
	QoS2
)

// QoS represents the Quality of Service levels.
type QoS uint8

// Valid returns true if the Qos value is valid.
func (q QoS) Valid() bool {
	return q == QoS0 || q == QoS1 || q == QoS2
}
