package packet

import (
	"testing"

	"github.com/yosssi/gmq/mqtt"
)

func TestCONNECT_setFixedHeader(t *testing.T) {
	p := &CONNECT{}

	p.setFixedHeader()

	want := []byte{TypeCONNECT << 4, 0x00}

	if len(p.fixedHeader) != len(want) || p.fixedHeader[0] != want[0] || p.fixedHeader[1] != want[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestCONNECT_setVariableHeader(t *testing.T) {
	p := &CONNECT{}

	p.setVariableHeader()

	want := []byte{0x00, 0x04, 0x4D, 0x51, 0x54, 0x54, 0x04, 0x00, 0x00, 0x00}

	if len(p.variableHeader) != len(want) {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
		return
	}

	for i, b := range p.variableHeader {
		if b != want[i] {
			t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
			return
		}
	}
}

func TestCONNECT_setPayload(t *testing.T) {
	p := &CONNECT{
		clientID:    []byte("clientID"),
		willTopic:   []byte("willTopic"),
		willMessage: []byte("willMessage"),
		userName:    []byte("userName"),
		password:    []byte("password"),
	}

	p.setPayload()

	want := []byte{0x00, 0x08, 0x63, 0x6C, 0x69, 0x65, 0x6E, 0x74, 0x49, 0x44, 0x00, 0x09, 0x77, 0x69, 0x6C, 0x6C, 0x54, 0x6F, 0x70, 0x69, 0x63, 0x00, 0x0B, 0x77, 0x69, 0x6C, 0x6C, 0x4D, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x00, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4E, 0x61, 0x6D, 0x65, 0x00, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6F, 0x72, 0x64}

	if len(p.payload) != len(want) {
		t.Errorf("p.payload=> %v, want => %v", p.payload, want)
		return
	}

	for i, b := range p.payload {
		if b != want[i] {
			t.Errorf("p.payload => %v, want => %v", p.payload, want)
			return
		}
	}
}

func TestCONNECT_connectFlags(t *testing.T) {
	p := &CONNECT{
		userName:     []byte("userName"),
		password:     []byte("password"),
		willRetain:   true,
		willQoS:      mqtt.QoS2,
		willTopic:    []byte("willTopic"),
		willMessage:  []byte("willMessage"),
		cleanSession: true,
	}

	b := p.connectFlags()

	if want := byte(0xF6); b != want {
		t.Errorf("b => %X, want => %X", b, want)
	}
}

func TestCONNECT_will(t *testing.T) {
	testCases := []struct {
		in  *CONNECT
		out bool
	}{
		{in: &CONNECT{}, out: false},
		{in: &CONNECT{willTopic: []byte{0x00}}, out: false},
		{in: &CONNECT{willMessage: []byte{0x00}}, out: false},
		{in: &CONNECT{willTopic: []byte{0x00}, willMessage: []byte{0x00}}, out: true},
	}

	for _, tc := range testCases {
		if got := tc.in.will(); got != tc.out {
			t.Errorf("got => %t, want => %t", got, tc.out)
		}
	}
}

func TestNewCONNECT_optsNil(t *testing.T) {
	if _, err := NewCONNECT(nil); err != ErrInvalidClientIDCleanSession {
		invalidError(t, err, ErrInvalidClientIDCleanSession)
	}
}

func TestNewCONNECT(t *testing.T) {
	p, err := NewCONNECT(&CONNECTOptions{
		ClientID: []byte("clientID"),
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

	if ptype != TypeCONNECT {
		t.Errorf("ptype => %X, want => %X", ptype, TypeCONNECT)
	}
}
