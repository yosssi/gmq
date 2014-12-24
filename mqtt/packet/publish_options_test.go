package packet

import "testing"

func TestPUBLISHOptions_validate_ErrInvalidQoS(t *testing.T) {
	opts := &PUBLISHOptions{
		QoS: 0x03,
	}

	if err := opts.validate(); err != ErrInvalidQoS {
		invalidError(t, err, ErrInvalidQoS)
	}
}

func TestPUBLISHOptions_validate_ErrTopicNameExceedsMaxStringsLen(t *testing.T) {
	opts := &PUBLISHOptions{
		TopicName: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrTopicNameExceedsMaxStringsLen {
		invalidError(t, err, ErrTopicNameExceedsMaxStringsLen)
	}
}

func TestPUBLISHOptions_validate_ErrTopicNameContainsWildcards(t *testing.T) {
	sliceOpts := []*PUBLISHOptions{
		&PUBLISHOptions{
			TopicName: []byte(wildcardMulti),
		},
		&PUBLISHOptions{
			TopicName: []byte(wildcardSingle),
		},
		&PUBLISHOptions{
			TopicName: []byte(wildcards),
		},
	}

	for _, opts := range sliceOpts {
		if err := opts.validate(); err != ErrTopicNameContainsWildcards {
			invalidError(t, err, ErrTopicNameContainsWildcards)
		}
	}
}

func TestPUBLISHOptions_validate_ErrMessageExceedsMaxStringsLen(t *testing.T) {
	opts := &PUBLISHOptions{
		Message: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrMessageExceedsMaxStringsLen {
		invalidError(t, err, ErrMessageExceedsMaxStringsLen)
	}
}

func TestPUBLISHOptions_validate(t *testing.T) {
	opts := &PUBLISHOptions{}

	if err := opts.validate(); err != nil {
		nilErrorExpected(t, err)
	}
}
