package packet

import "testing"

func Test_pubrel_PacketID(t *testing.T) {
	want := uint16(1)

	p := &pubrel{
		packetID: want,
	}

	if got := p.PacketID(); got != want {
		t.Errorf("got => %d, want => %d", got, want)
	}
}

func Test_pubrel_setFixedHeader(t *testing.T) {
	p := &pubrel{
		packetID: 1,
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

func Test_pubrel_setVariableHeader(t *testing.T) {
	p := &pubrel{
		packetID: 65534,
	}

	p.setVariableHeader()

	want := []byte{0xFF, 0xFE}

	if len(p.variableHeader) != len(want) || p.variableHeader[0] != want[0] || p.variableHeader[1] != want[1] {
		t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
	}
}

func TestNewPUBREL(t *testing.T) {
	p := NewPUBREL(nil)

	ptype, err := p.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypePUBREL {
		t.Errorf("ptype => %X, want => %X", ptype, TypePUBREL)
	}
}
