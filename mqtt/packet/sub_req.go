package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Error value
var ErrTopicFilterExceedsMaxStringsLen = errors.New("the length of the Topic Filter exceeds the maximum strings legnth")

// SubReq represents subscription request.
type SubReq struct {
	// TopicFilter is the Topic Filter of the Subscription.
	TopicFilter []byte
	// QoS is the requsting QoS.
	QoS byte
}

// validate validates the subscription request.
func (s *SubReq) validate() error {
	// Check the length of the Topic Filter.
	if len(s.TopicFilter) > maxStringsLen {
		return ErrTopicFilterExceedsMaxStringsLen
	}

	// Check the QoS.
	if !mqtt.ValidQoS(s.QoS) {
		return ErrInvalidQoS
	}

	return nil
}
