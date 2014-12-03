package packet

import "testing"

func Test_encodeUint16(t *testing.T) {
	testCases := []struct {
		in  uint16
		out []byte
	}{
		{
			0x0000,
			[]byte{0x00, 0x00},
		},
		{
			0x0001,
			[]byte{0x00, 0x01},
		},
		{
			0x00FF,
			[]byte{0x00, 0xFF},
		},
		{
			0x0100,
			[]byte{0x01, 0x00},
		},
		{
			0xFF00,
			[]byte{0xFF, 0x00},
		},
		{
			0x0101,
			[]byte{0x01, 0x01},
		},
		{
			0x01FF,
			[]byte{0x01, 0xFF},
		},
		{
			0xFF01,
			[]byte{0xFF, 0x01},
		},
		{
			0xFFFF,
			[]byte{0xFF, 0xFF},
		},
	}

	for _, tc := range testCases {
		result := encodeUint16(tc.in)

		if len(result) != len(tc.out) {
			t.Errorf("encodeUint16(%d) => %v, want %v", tc.in, result, tc.out)
			continue
		}

		for i, b := range result {
			if b != tc.out[i] {
				t.Errorf("encodeUint16(%d) => %v, want %v", tc.in, result, tc.out)
				break
			}
		}
	}
}

func Test_encodeLength(t *testing.T) {
	testCases := []struct {
		in  uint
		out uint32
	}{
		{127, 127},
		{16383, 65407},
		{2097151, 16777087},
		{268435455, 4294967167},
	}

	for _, tc := range testCases {
		if result := encodeLength(tc.in); result != tc.out {
			t.Errorf("encodeLength(%d) => %d, want %d", tc.in, result, tc.out)
		}
	}
}
