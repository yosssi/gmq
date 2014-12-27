package packet

import "testing"

func TestPUBACK_setFixedHeader(t *testing.T) {
	p := &PUBACK{
		PacketID: 1,
	}

	p.setFixedHeader()

	want := []byte{0x40, 0x00}

	if len(want) != len(p.fixedHeader) || want[0] != p.fixedHeader[0] || want[1] != p.fixedHeader[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestPUBACK_setVariableHeader(t *testing.T) {
	p := &PUBACK{
		PacketID: 1,
	}

	p.setVariableHeader()

	want := []byte{0x00, 0x01}

	if len(want) != len(p.variableHeader) || want[0] != p.variableHeader[0] || want[1] != p.variableHeader[1] {
		t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
	}
}

func TestNewPUBACK_optsNil(t *testing.T) {
	if _, err := NewPUBACK(nil); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestNewPUBACK(t *testing.T) {
	_, err := NewPUBACK(&PUBACKOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewPUBACKFromBytes_validatePUBACKBytesErr(t *testing.T) {
	if _, err := NewPUBACKFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBACKFromBytes(t *testing.T) {
	if _, err := NewPUBACKFromBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
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

func Test_validatePUBACKBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x00}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func Test_validatePUBACKBytes(t *testing.T) {
	if err := validatePUBACKBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}
