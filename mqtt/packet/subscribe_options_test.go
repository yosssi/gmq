package packet

import "testing"

func TestSUBSCRIBEOptions_validate_ErrInvalidNoSubReq(t *testing.T) {
	opts := &SUBSCRIBEOptions{
		PacketID: 1,
	}

	if err := opts.validate(); err != ErrInvalidNoSubReq {
		invalidError(t, err, ErrInvalidNoSubReq)
	}
}

func TestSUBSCRIBEOptions_validate_ErrNoTopicFilter(t *testing.T) {
	opts := &SUBSCRIBEOptions{
		PacketID: 1,
		SubReqs: []*SubReq{
			&SubReq{},
		},
	}

	if err := opts.validate(); err != ErrNoTopicFilter {
		invalidError(t, err, ErrNoTopicFilter)
	}
}
