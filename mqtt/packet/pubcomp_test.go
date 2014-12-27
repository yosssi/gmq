package packet

import "testing"

func TestPUBCOMP_setFixedHeader(t *testing.T) {
	p := &PUBCOMP{
		PacketID: 1,
	}

	p.setFixedHeader()

	want := []byte{0x70, 0x00}

	if len(want) != len(p.fixedHeader) || want[0] != p.fixedHeader[0] || want[1] != p.fixedHeader[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestPUBCOMP_setVariableHeader(t *testing.T) {
	p := &PUBCOMP{
		PacketID: 1,
	}

	p.setVariableHeader()

	want := []byte{0x00, 0x01}

	if len(want) != len(p.variableHeader) || want[0] != p.variableHeader[0] || want[1] != p.variableHeader[1] {
		t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
	}
}

func TestNewPUBCOMP_optsNil(t *testing.T) {
	if _, err := NewPUBCOMP(nil); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestNewPUBCOMP(t *testing.T) {
	_, err := NewPUBCOMP(&PUBCOMPOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewPUBCOMPFromBytes_validatePUBCOMPBytesErr(t *testing.T) {
	if _, err := NewPUBCOMPFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBCOMPFromBytes(t *testing.T) {
	if _, err := NewPUBCOMPFromBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
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

func Test_validatePUBCOMPBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x00}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func Test_validatePUBCOMPBytes(t *testing.T) {
	if err := validatePUBCOMPBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}
