package packet

import "testing"

func Test_encodeUint16(t *testing.T) {
	b := encodeUint16(0x0123)

	want := []byte{0x01, 0x23}

	if len(b) != len(want) || b[0] != want[0] || b[1] != want[1] {
		t.Errorf("b => %v, want => %v", b, want)
	}
}

func Test_encodeLength(t *testing.T) {
	testCases := []struct {
		in  uint32
		out uint32
	}{
		{in: 0, out: 0x00},
		{in: 127, out: 0x7F},
		{in: 128, out: 0x8001},
		{in: 16383, out: 0xFF7F},
		{in: 16384, out: 0x808001},
		{in: 2097151, out: 0xFFFF7F},
		{in: 2097152, out: 0x80808001},
		{in: 268435455, out: 0xFFFFFF7F},
	}

	for _, tc := range testCases {
		if got := encodeLength(tc.in); got != tc.out {
			t.Errorf("got => %d, want => %d", got, tc.out)
		}
	}
}
