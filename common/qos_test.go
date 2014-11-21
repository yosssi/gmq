package common

import "testing"

func TestQoS_Valid(t *testing.T) {
	qos := QoS(0)
	qos.Valid()
}
