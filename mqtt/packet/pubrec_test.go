package packet

import "testing"

func TestNewPUBRECFromBytes_validatePUBRECBytesErr(t *testing.T) {
	if _, err := NewPUBRECFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBRECFromBytes(t *testing.T) {
	if _, err := NewPUBRECFromBytes([]byte{TypePUBREC << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_validatePUBRECBytes_ptypeErr(t *testing.T) {
	if err := validatePUBRECBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBRECBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypePUBREC << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBRECBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePUBRECBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypePUBREC<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validatePUBRECBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypePUBREC << 4, 0x00}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePUBRECBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypePUBREC << 4, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validatePUBRECBytes(t *testing.T) {
	if err := validatePUBRECBytes([]byte{TypePUBREC << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}
