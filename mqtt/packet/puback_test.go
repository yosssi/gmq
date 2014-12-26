package packet

import "testing"

func TestNewPUBACKFromBytes_validatePUBACKBytesErr(t *testing.T) {
	if _, err := NewPUBACKFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBACKFromBytes(t *testing.T) {
	if _, err := NewPUBACKFromBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_validatePUBACKBytes_ptypeErr(t *testing.T) {
	if err := validatePUBACKBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBACKBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBACKBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePUBACKBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validatePUBACKBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4, 0x00}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePUBACKBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validatePUBACKBytes(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}
