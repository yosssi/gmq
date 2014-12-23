package packet

import "testing"

func TestNewPINGRESPFromBytes_errValidatePINGRESPBytes(t *testing.T) {
	if _, err := NewPINGRESPFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPINGRESPFromBytes(t *testing.T) {
	if _, err := NewPINGRESPFromBytes([]byte{TypePINGRESP << 4, 0x00}, nil); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_validatePINGRESPBytes_ptypeErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePINGRESPBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePINGRESPBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePINGRESPBytes([]byte{TypePINGRESP << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePINGRESPBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePINGRESPBytes([]byte{0x00 << 4, 0x00}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePINGRESPBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validatePINGRESPBytes([]byte{TypePINGRESP<<4 | 0x01, 0x00}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validatePINGRESPBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePINGRESPBytes([]byte{TypePINGRESP << 4, 0x01}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePINGRESPBytes_ErrInvalidRemainingLen(t *testing.T) {
	if err := validatePINGRESPBytes([]byte{TypePINGRESP << 4, 0x00}, []byte{0x00}); err != ErrInvalidRemainingLen {
		invalidError(t, err, ErrInvalidRemainingLen)
	}
}
