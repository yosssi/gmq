package packet

import (
	"testing"

	"github.com/yosssi/gmq/mqtt"
)

func TestPUBLISH_setFixedHeader(t *testing.T) {
	p := &PUBLISH{
		DUP:    true,
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
		QoS:       mqtt.QoS1,
		TopicName: []byte("topicName"),
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
		Message: []byte{0x00, 0x01},
	}

	p.setPayload()

	if len(p.payload) != len(p.Message) || p.payload[0] != p.Message[0] || p.payload[1] != p.Message[1] {
		t.Errorf("p.payload => %v, want => %v", p.payload, p.Message)
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

func TestNewPUBLISHFromBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if _, err := NewPUBLISHFromBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewPUBLISHFromBytes(t *testing.T) {
	if _, err := NewPUBLISHFromBytes([]byte{TypePUBLISH<<4 | 0x02, 0x07}, []byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_validatePUBLISHBytes_fixedHeaderErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBLISHBytes(nil, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidFixedHeaderLen(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH << 4}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidPacketType(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypeCONNECT << 4, 0x00}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidQoS(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH<<4 | 0x06, 0x00}, nil); err != ErrInvalidQoS {
		invalidError(t, err, ErrInvalidQoS)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidRemainingLen(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH << 4, 0x00}, nil); err != ErrInvalidRemainingLen {
		invalidError(t, err, ErrInvalidRemainingLen)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidRemainingLength(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH<<4 | 0x02, 0x02}, []byte{0x00, 0x00}); err != ErrInvalidRemainingLength {
		invalidError(t, err, ErrInvalidRemainingLength)
	}
}

func Test_validatePUBLISHBytes_ErrInvalidPacketID(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH<<4 | 0x02, 0x07}, []byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x00}); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func Test_validatePUBLISHBytes(t *testing.T) {
	if err := validatePUBLISHBytes([]byte{TypePUBLISH<<4 | 0x02, 0x07}, []byte{0x00, 0x03, 0x61, 0x2F, 0x62, 0x00, 0x01}); err != nil {
		nilErrorExpected(t, err)
	}
}
