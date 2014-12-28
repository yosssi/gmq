package packet

import "testing"

func TestUNSUBSCRIBE_setFixedHeader(t *testing.T) {
	p := &UNSUBSCRIBE{}

	p.variableHeader = make([]byte, 1)

	p.payload = make([]byte, 2)

	p.setFixedHeader()

	want := []byte{0xA2, 0x03}

	if len(p.fixedHeader) != len(want) || p.fixedHeader[0] != want[0] || p.fixedHeader[1] != want[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestUNSUBSCRIBE_setVariableHeader(t *testing.T) {
	p := &UNSUBSCRIBE{
		PacketID: 1,
	}

	p.setVariableHeader()

	want := []byte{0x00, 0x01}

	if len(p.variableHeader) != len(want) || p.variableHeader[0] != want[0] || p.variableHeader[1] != want[1] {
		t.Errorf("p.variableHeader=> %v, want => %v", p.variableHeader, want)
	}
}

func TestUNSUBSCRIBE_setPayload(t *testing.T) {
	p := &UNSUBSCRIBE{
		TopicFilters: [][]byte{
			[]byte("t"),
		},
	}

	p.setPayload()

	want := []byte{0x00, 0x01, 0x74}

	if len(p.payload) != len(want) || p.payload[0] != want[0] || p.payload[1] != want[1] {
		t.Errorf("p.payload => %v, want => %v", p.payload, want)
	}
}

func TestNewUNSUBSCRIBE_optsNil(t *testing.T) {
	if _, err := NewUNSUBSCRIBE(nil); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestNewUNSUBSCRIBE(t *testing.T) {
	_, err := NewUNSUBSCRIBE(&UNSUBSCRIBEOptions{
		PacketID: 1,
		TopicFilters: [][]byte{
			[]byte("t"),
		},
	})
	if err != nil {
		nilErrorExpected(t, err)
	}
}
