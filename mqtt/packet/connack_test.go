package packet

import "testing"

func TestNewCONNACKFromBytes_errValidateCONNACKBytes(t *testing.T) {
	if _, err := NewCONNACKFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewCONNACKFromBytes(t *testing.T) {
	p, err := NewCONNACKFromBytes([]byte{TypeCONNACK << 4, 0x02}, []byte{0x00, 0x00})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	ptype, err := p.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypeCONNACK {
		t.Errorf("ptype => %X, want => %X", ptype, TypeCONNACK)
	}
}

func Test_validateCONNACKBytes_ptypeErr(t *testing.T) {
	if err := validateCONNACKBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateCONNACKBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validateCONNACKBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNECT << 4, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validateCONNACKBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validateCONNACKBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4, 0x01}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validateCONNACKBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validateCONNACKBytes_ErrInvalidVariableHeader(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4, 0x02}, []byte{0x02, 0x00}); err != ErrInvalidVariableHeader {
		invalidError(t, err, ErrInvalidVariableHeader)
	}
}

func Test_validateCONNACKBytes_ErrInvalidConnectReturnCode(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4, 0x02}, []byte{0x00, 0x06}); err != ErrInvalidConnectReturnCode {
		invalidError(t, err, ErrInvalidConnectReturnCode)
	}
}

func Test_validateCONNACKBytes(t *testing.T) {
	if err := validateCONNACKBytes([]byte{TypeCONNACK << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}
