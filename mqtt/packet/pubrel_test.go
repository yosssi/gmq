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
