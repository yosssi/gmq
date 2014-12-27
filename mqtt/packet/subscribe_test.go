package packet

import "testing"

func TestSUBSCRIBE_setFixedHeader(t *testing.T) {
	p := &SUBSCRIBE{}

	p.variableHeader = make([]byte, 2)
	p.payload = make([]byte, 4)

	p.setFixedHeader()

	want := []byte{0x82, 0x06}

	if len(p.fixedHeader) != len(want) || p.fixedHeader[0] != want[0] || p.fixedHeader[1] != want[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestSUBSCRIBE_setVariableHeader(t *testing.T) {
	p := &SUBSCRIBE{
		PacketID: 1,
	}

	p.setVariableHeader()

	want := []byte{0x00, 0x01}

	if len(p.variableHeader) != len(want) || p.variableHeader[0] != want[0] || p.variableHeader[1] != want[1] {
		t.Errorf("p.variableHeader=> %v, want => %v", p.variableHeader, want)
	}
}

func TestSUBSCRIBE_setPayload(t *testing.T) {
	p := &SUBSCRIBE{
		SubReqs: []*SubReq{
			&SubReq{
				TopicFilter: []byte("t"),
			},
		},
	}

	p.setPayload()

	want := []byte{0x00, 0x01, 0x74, 0x00}

	if len(p.payload) != len(want) || p.payload[0] != want[0] || p.payload[1] != want[1] || p.payload[2] != want[2] || p.payload[3] != want[3] {
		t.Errorf("p.payload=> %v, want => %v", p.payload, want)
	}
}

func TestNewSUBSCRIBE_optsNil(t *testing.T) {
	if _, err := NewSUBSCRIBE(nil); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestNewSUBSCRIBE(t *testing.T) {
	_, err := NewSUBSCRIBE(&SUBSCRIBEOptions{
		PacketID: 1,
		SubReqs: []*SubReq{
			&SubReq{
				TopicFilter: []byte("t"),
			},
		},
	})
	if err != nil {
		nilErrorExpected(t, err)
	}
}
