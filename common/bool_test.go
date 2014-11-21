package common

import "testing"

func TestBoolPtr(t *testing.T) {
	testCases := []bool{false, true}

	for _, tc := range testCases {
		if b := BoolPtr(tc); *b != tc {
			t.Errorf("BoolPtr(%t) => %t, want %t", tc, *b, tc)
		}
	}
}
