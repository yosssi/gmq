package packet

import "testing"

func TestUNSUBSCRIBEOptions_validate_ErrInvalidPacketID(t *testing.T) {
	opts := &UNSUBSCRIBEOptions{}

	if err := opts.validate(); err != ErrInvalidPacketID {
		invalidError(t, err, ErrInvalidPacketID)
	}
}

func TestUNSUBSCRIBEOptions_validate_ErrNoTopicFilter(t *testing.T) {
	opts := &UNSUBSCRIBEOptions{
		PacketID: 1,
	}

	if err := opts.validate(); err != ErrNoTopicFilter {
		invalidError(t, err, ErrNoTopicFilter)
	}
}

func TestUNSUBSCRIBEOptions_validate_topicFilterErrNoTopicFilter(t *testing.T) {
	opts := &UNSUBSCRIBEOptions{
		PacketID:     1,
		TopicFilters: make([][]byte, 1),
	}

	if err := opts.validate(); err != ErrNoTopicFilter {
		invalidError(t, err, ErrNoTopicFilter)
	}
}

func TestUNSUBSCRIBEOptions_validate_ErrTopicFilterExceedsMaxStringsLen(t *testing.T) {
	opts := &UNSUBSCRIBEOptions{
		PacketID: 1,
		TopicFilters: [][]byte{
			make([]byte, maxStringsLen+1),
		},
	}

	if err := opts.validate(); err != ErrTopicFilterExceedsMaxStringsLen {
		invalidError(t, err, ErrTopicFilterExceedsMaxStringsLen)
	}
}

func TestUNSUBSCRIBEOptions_validate(t *testing.T) {
	opts := &UNSUBSCRIBEOptions{
		PacketID: 1,
		TopicFilters: [][]byte{
			make([]byte, 1),
		},
	}

	if err := opts.validate(); err != nil {
		nilErrorExpected(t, err)
	}
}
