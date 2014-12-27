package packet

import "testing"

func TestNewSUBACKFromBytes_err(t *testing.T) {
	if _, err := NewSUBACKFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateSUBACKBytes_ptypeErr(t *testing.T) {
	if err := validateSUBACKBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateSUBACKBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeSUBACK << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateSUBACKBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validateSUBACKBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeSUBACK<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validateSUBACKBytes_ErrInvalidRemainingLen(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeSUBACK << 4, 0x02}, []byte{0x00, 0x01}); err != ErrInvalidRemainingLen {
		invalidError(t, err, ErrInvalidRemainingLen)
	}
}

func Test_validateSUBACKBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeSUBACK << 4, 0x03}, []byte{0x00, 0x00, 0x02}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func Test_validateSUBACKBytes_ErrInvalidSUBACKReturnCode(t *testing.T) {
	if err := validateSUBACKBytes([]byte{TypeSUBACK << 4, 0x03}, []byte{0x00, 0x01, 0x03}); err != ErrInvalidSUBACKReturnCode {
		invalidError(t, err, ErrInvalidSUBACKReturnCode)
	}
}
