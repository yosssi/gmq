package packet

import "testing"

func TestNewFromBytes_errTypeFromBytes(t *testing.T) {
	if _, err := NewFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrInvalidFixedHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrInvalidFixedHeaderLen)
		}
	}
}

func TestNewFromBytes_CONNACK(t *testing.T) {
	fixedHeader := []byte{TypeCONNACK << 4, 0x02}
	remaining := []byte{0x00, 0x00}

	if _, err := NewFromBytes(fixedHeader, remaining); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNewFromBytes_errNewCONNACKFromBytes(t *testing.T) {
	fixedHeader := []byte{TypeCONNACK << 4, 0x00}
	remaining := []byte{0x00, 0x00}

	if _, err := NewFromBytes(fixedHeader, remaining); err != ErrCONNACKInvalidRemainingLength {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidRemainingLength)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidRemainingLength)
		}
	}
}

func TestNewFromBytes_PINGRESP(t *testing.T) {
	fixedHeader := []byte{TypePINGRESP << 4, 0x00}

	if _, err := NewFromBytes(fixedHeader, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNewFromBytes_errNewPINGRESPFromBytes(t *testing.T) {
	fixedHeader := []byte{TypePINGRESP << 4}

	if _, err := NewFromBytes(fixedHeader, nil); err != ErrPINGRESPInvalidFixedHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrPINGRESPInvalidFixedHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrPINGRESPInvalidFixedHeaderLen)
		}
	}
}

func TestNewFromBytes_errInvalidPacketType(t *testing.T) {
	fixedHeader := []byte{0x00}

	if _, err := NewFromBytes(fixedHeader, nil); err != ErrInvalidPacketType {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrInvalidPacketType)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrInvalidPacketType)
		}
	}
}
