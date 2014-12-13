package packet

import "testing"

func TestNewPINGRESPFromBytes_errPINGRESPInvalidFixedHeaderLen(t *testing.T) {
	if _, err := NewPINGRESPFromBytes(nil); err != ErrPINGRESPInvalidFixedHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrPINGRESPInvalidFixedHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrPINGRESPInvalidFixedHeaderLen)
		}
	}
}

func TestNewPINGRESPFromBytes_errPINGRESPInvalidFixedHeaderFirstByte(t *testing.T) {
	if _, err := NewPINGRESPFromBytes([]byte{0x00, 0x00}); err != ErrPINGRESPInvalidFixedHeaderFirstByte {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrPINGRESPInvalidFixedHeaderFirstByte)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrPINGRESPInvalidFixedHeaderFirstByte)
		}
	}
}

func TestNewPINGRESPFromBytes_errPINGRESPInvalidRemainingLength(t *testing.T) {
	if _, err := NewPINGRESPFromBytes([]byte{TypePINGRESP << 4, 0x01}); err != ErrPINGRESPInvalidRemainingLength {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrPINGRESPInvalidRemainingLength)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrPINGRESPInvalidRemainingLength)
		}
	}
}

func TestNewPINGRESPFromBytes(t *testing.T) {
	if _, err := NewPINGRESPFromBytes([]byte{TypePINGRESP << 4, 0x00}); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}
