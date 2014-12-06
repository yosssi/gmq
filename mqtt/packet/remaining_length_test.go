package packet

import "testing"

func Test_appendRemainingLength(t *testing.T) {
	testCase := []struct {
		in  uint32
		out []byte
	}{
		{
			in:  0xFF000000,
			out: []byte{255, 0, 0, 0},
		},
		{
			in:  0x00FF0000,
			out: []byte{255, 0, 0},
		},
		{
			in:  0x0000FF00,
			out: []byte{255, 0},
		},
		{
			in:  0x000000FF,
			out: []byte{255},
		},
	}

	for _, tc := range testCase {
		b := appendRemainingLength([]byte{}, tc.in)

		if len(b) != len(tc.out) {
			t.Errorf("len(b) => %d, want => %d", len(b), len(tc.out))
			continue
		}

		for i, bt := range b {
			if bt != tc.out[i] {
				t.Errorf("bt => %X, want => %X", bt, tc.out[i])
				continue
			}
		}
	}
}
