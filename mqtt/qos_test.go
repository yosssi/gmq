package mqtt

import "testing"

func TestValidQoS(t *testing.T) {
	testCases := []struct {
		in  uint
		out bool
	}{
		{in: QoS0, out: true},
		{in: QoS1, out: true},
		{in: QoS2, out: true},
		{in: 3, out: false},
	}

	for _, tc := range testCases {
		if got := ValidQoS(tc.in); got != tc.out {
			t.Errorf("ValidQoS(tc.in) => %t, want => %t", got, tc.out)
		}
	}
}
