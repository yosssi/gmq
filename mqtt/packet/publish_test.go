package packet

import "testing"

func TestPUBLISH_setFixedHeader(t *testing.T) {
	p := &PUBLISH{}

	p.DUP = true
	p.Retain = true
	p.TopicName = "topicName"
	p.Message = "message"

	p.setVariableHeader()

	p.setPayload()

	p.setFixedHeader()

	if got, want := p.FixedHeader[0], byte(0xB1); got != want {
		t.Errorf("got => %d, want => %d", got, want)
	}

	if got, want := int(p.FixedHeader[1]), len(p.VariableHeader)+len(p.Payload); got != want {
		t.Errorf("got => %d, want => %d", got, want)
	}
}

func TestNewPUBLISH_optsNil(t *testing.T) {
	if _, err := NewPUBLISH(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNewPUBLISH_errPUBLISHInvalidQoS(t *testing.T) {
	if _, err := NewPUBLISH(&PUBLISHOptions{QoS: 3}); err != ErrPUBLISHInvalidQoS {
		errorfErr(t, err, ErrPUBLISHInvalidQoS)
	}
}

func TestNewPUBLISH_errPUBLISHTopicNameContainsWildcard(t *testing.T) {
	if _, err := NewPUBLISH(&PUBLISHOptions{TopicName: "test*test"}); err != ErrPUBLISHTopicNameContainsWildcard {
		errorfErr(t, err, ErrPUBLISHTopicNameContainsWildcard)
	}
}
