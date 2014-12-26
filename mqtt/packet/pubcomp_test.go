package packet

import "testing"

func TestNewPUBCOMPFromBytes_validatePUBCOMPBytesErr(t *testing.T) {
	if _, err := NewPUBCOMPFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBCOMPFromBytes(t *testing.T) {
	if _, err := NewPUBCOMPFromBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_validatePUBCOMPBytes_ptypeErr(t *testing.T) {
	if err := validatePUBCOMPBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBCOMPBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBCOMPBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePUBCOMPBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validatePUBCOMPBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4, 0x00}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePUBCOMPBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validatePUBCOMPBytes(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}
