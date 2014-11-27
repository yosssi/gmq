package packet

import "testing"

func TestNewCONNACKFromBytes_errInvalidCONNACKVariableHeaderLen(t *testing.T) {
	if _, err := NewCONNACKFromBytes(nil, nil); err != ErrInvalidCONNACKVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrInvalidCONNACKVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrInvalidCONNACKVariableHeaderLen)
		}
	}
}

func TestNewCONNACKFromBytes(t *testing.T) {
	if _, err := NewCONNACKFromBytes(nil, []byte{0, 0}); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}
