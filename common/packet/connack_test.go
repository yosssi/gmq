package packet

import "testing"

func TestNewCONNACKFromBytes_errInvalidCONNACKVariableHeaderLen(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x02}, nil); err != ErrInvalidCONNACKVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrInvalidCONNACKVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrInvalidCONNACKVariableHeaderLen)
		}
	}
}

func TestNewCONNACKFromBytes(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x02}, []byte{0x00, 0x00}); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}
