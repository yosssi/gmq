package packet

import "testing"

func Test_decodeUint16_ErrInvalidByteLen(t *testing.T) {
	if _, err := decodeUint16(nil); err != ErrInvalidByteLen {
		invalidError(t, err, ErrInvalidByteLen)
	}
}

func Test_decodeUint16(t *testing.T) {
	var i uint16

	for {
		b := encodeUint16(i)

		u, err := decodeUint16(b)
		if err != nil {
			nilErrorExpected(t, err)
			continue
		}

		if u != i {
			t.Errorf("u => %d, want => %d", u, i)
		}

		if i == 65535 {
			break
		}

		i++
	}
}
