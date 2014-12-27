package packet

import "testing"

func TestSubReq_validate_ErrNoTopicFilter(t *testing.T) {
	s := &SubReq{}

	if err := s.validate(); err != ErrNoTopicFilter {
		invalidError(t, err, ErrNoTopicFilter)
	}
}

func TestSubReq_validate_ErrTopicFilterExceedsMaxStringsLen(t *testing.T) {
	s := &SubReq{
		TopicFilter: make([]byte, maxStringsLen+1),
	}

	if err := s.validate(); err != ErrTopicFilterExceedsMaxStringsLen {
		invalidError(t, err, ErrTopicFilterExceedsMaxStringsLen)
	}
}

func TestSubReq_validate_ErrInvalidQoS(t *testing.T) {
	s := &SubReq{
		TopicFilter: make([]byte, 1),
		QoS:         0x03,
	}

	if err := s.validate(); err != ErrInvalidQoS {
		invalidError(t, err, ErrInvalidQoS)
	}
}

func TestSubReq_validate(t *testing.T) {
	s := &SubReq{
		TopicFilter: make([]byte, 1),
	}

	if err := s.validate(); err != nil {
		nilErrorExpected(t, err)
	}
}
