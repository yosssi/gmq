package packet

import (
	"testing"

	"github.com/yosssi/gmq/mqtt"
)

func TestPUBLISH_setFixedHeader(t *testing.T) {
	p := &PUBLISH{
		dup:    true,
		retain: true,
	}

	p.variableHeader = []byte{0x00}

	p.payload = []byte{0x00, 0x00}

	p.setFixedHeader()

	want := []byte{0x39, 0x03}

	if len(p.fixedHeader) != len(want) || p.fixedHeader[0] != want[0] || p.fixedHeader[1] != want[1] {
		t.Errorf("p.fixedHeader => %v, want => %v", p.fixedHeader, want)
	}
}

func TestPUBLISH_setVariableHeader(t *testing.T) {
	p := &PUBLISH{
		qos:       mqtt.QoS1,
		topicName: []byte("topicName"),
		PacketID:  1,
	}

	p.setVariableHeader()

	want := []byte{0x00, 0x09, 0x74, 0x6F, 0x70, 0x69, 0x63, 0x4E, 0x061, 0x6D, 0x65, 0x00, 0x01}

	if len(p.variableHeader) != len(want) {
		t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
		return
	}

	for i, b := range p.variableHeader {
		if b != want[i] {
			t.Errorf("p.variableHeader => %v, want => %v", p.variableHeader, want)
			return
		}
	}
}

func TestPUBLISH_setPayload(t *testing.T) {
	p := &PUBLISH{
		message: []byte{0x00, 0x01},
	}

	p.setPayload()

	if len(p.payload) != len(p.message) || p.payload[0] != p.message[0] || p.payload[1] != p.message[1] {
		t.Errorf("p.payload => %v, want => %v", p.payload, p.message)
	}
}

func TestNewPUBLISH_optsNil(t *testing.T) {
	if _, err := NewPUBLISH(nil); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewPUBLISH_validateErr(t *testing.T) {
	_, err := NewPUBLISH(&PUBLISHOptions{
		QoS: 0x03,
	})

	if err != ErrInvalidQoS {
		invalidError(t, err, ErrInvalidQoS)
	}
}
