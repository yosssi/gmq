package packet

import "testing"

func TestNewUNSUBACKFromBytes_err(t *testing.T) {
	if _, err := NewUNSUBACKFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateUNSUBACKBytes_ptypeErr(t *testing.T) {
	if err := validateUNSUBACKBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeUNSUBACK << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeUNSUBACK<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeUNSUBACK << 4, 0x00}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeUNSUBACK << 4, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validateUNSUBACKBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validateUNSUBACKBytes([]byte{TypeUNSUBACK << 4, 0x02}, []byte{0x00, 0x00}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}
