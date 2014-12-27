package packet

import "testing"

func TestPUBREL_setFixedHeader(t *testing.T) {
	p := &PUBREL{
		PacketID: 1,
	}

	p.variableHeader = []byte{0x00, 0x01}

	p.setFixedHeader()

	want := []byte{
		TypePUBREL<<4 | 0x02,
		byte(len(p.variableHeader)),
	}

	if len(p.fixedHeader) != len(want) || p.fixedHeader[0] != want[0] || p.fixedHeader[1] != want[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestPUBREL_setVariableHeader(t *testing.T) {
	p := &PUBREL{
		PacketID: 65534,
	}

	p.setVariableHeader()

	want := []byte{0xFF, 0xFE}

	if len(p.variableHeader) != len(want) || p.variableHeader[0] != want[0] || p.variableHeader[1] != want[1] {
		t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
	}
}

func TestNewPUBREL_optsNil(t *testing.T) {
	if _, err := NewPUBREL(nil); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestNewPUBREL(t *testing.T) {
	p, err := NewPUBREL(&PUBRELOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	ptype, err := p.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypePUBREL {
		t.Errorf("ptype => %X, want => %X", ptype, TypePUBREL)
	}
}

func TestNewPUBRELFromBytes_validateErr(t *testing.T) {
	if _, err := NewPUBRELFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBRELBytes_ptypeErr(t *testing.T) {
	if err := validatePUBRELBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBRELBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypePUBREL<<4 | 0x02}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBRELBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypeCONNECT<<4 | 0x02, 0x02}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePUBRELBytes_ErrInvalidFixedHeader(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypePUBREL<<4 | 0x01, 0x02}, nil); err != ErrInvalidFixedHeader {
		invalidError(t, err, ErrInvalidFixedHeader)
	}
}

func Test_validatePUBRELBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypePUBREL<<4 | 0x02, 0x00}, nil); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePUBRELBytes_ErrInvalidVariableHeaderLen(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypePUBREL<<4 | 0x02, 0x02}, nil); err != ErrInvalidVariableHeaderLen {
		invalidError(t, err, ErrInvalidVariableHeaderLen)
	}
}

func Test_validatePUBRELBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validatePUBRELBytes([]byte{TypePUBREL<<4 | 0x02, 0x02}, []byte{0x00, 0x00}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}
