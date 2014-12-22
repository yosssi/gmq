package packet

import (
	"testing"

	"github.com/yosssi/gmq/mqtt"
)

func TestCONNECTOptions_validate_errClientIDExceedsMaxStringsLen(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrClientIDExceedsMaxStringsLen {
		invalidError(t, err, ErrClientIDExceedsMaxStringsLen)
	}
}

func TestCONNECTOptions_validate_errInvalidClientIDCleanSession(t *testing.T) {
	opts := &CONNECTOptions{}

	if err := opts.validate(); err != ErrInvalidClientIDCleanSession {
		invalidError(t, err, ErrInvalidClientIDCleanSession)
	}
}

func TestCONNECTOptions_validate_errUserNameExceedsMaxStringsLen(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
		UserName: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrUserNameExceedsMaxStringsLen {
		invalidError(t, err, ErrUserNameExceedsMaxStringsLen)
	}
}

func TestCONNECTOptions_validate_errPasswordExceedsMaxStringsLen(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
		Password: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrPasswordExceedsMaxStringsLen {
		invalidError(t, err, ErrPasswordExceedsMaxStringsLen)
	}
}

func TestCONNECTOptions_validate_errInvalidClientIDPassword(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
		Password: []byte{0x00},
	}

	if err := opts.validate(); err != ErrInvalidClientIDPassword {
		invalidError(t, err, ErrInvalidClientIDPassword)
	}
}

func TestCONNECTOptions_validate_errWillTopicExceedsMaxStringsLen(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID:  []byte{0x00},
		WillTopic: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrWillTopicExceedsMaxStringsLen {
		invalidError(t, err, ErrWillTopicExceedsMaxStringsLen)
	}
}

func TestCONNECTOptions_validate_errWillMessageExceedsMaxStringsLen(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID:    []byte{0x00},
		WillMessage: make([]byte, maxStringsLen+1),
	}

	if err := opts.validate(); err != ErrWillMessageExceedsMaxStringsLen {
		invalidError(t, err, ErrWillMessageExceedsMaxStringsLen)
	}
}

func TestCONNECTOptions_validate_errInvalidWillTopicMessage(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID:  []byte{0x00},
		WillTopic: []byte{0x00},
	}

	if err := opts.validate(); err != ErrInvalidWillTopicMessage {
		invalidError(t, err, ErrInvalidWillTopicMessage)
	}
}

func TestCONNECTOptions_validate_errInvalidWillQoS(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
		WillQoS:  byte(0x03),
	}

	if err := opts.validate(); err != ErrInvalidWillQoS {
		invalidError(t, err, ErrInvalidWillQoS)
	}
}

func TestCONNECTOptions_validate_errInvalidWillTopicMessageQoS(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
		WillQoS:  mqtt.QoS1,
	}

	if err := opts.validate(); err != ErrInvalidWillTopicMessageQoS {
		invalidError(t, err, ErrInvalidWillTopicMessageQoS)
	}
}

func TestCONNECTOptions_validate_errInvalidWillTopicMessageRetain(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID:   []byte{0x00},
		WillRetain: true,
	}

	if err := opts.validate(); err != ErrInvalidWillTopicMessageRetain {
		invalidError(t, err, ErrInvalidWillTopicMessageRetain)
	}
}

func TestCONNECTOptions_validate(t *testing.T) {
	opts := &CONNECTOptions{
		ClientID: []byte{0x00},
	}

	if err := opts.validate(); err != nil {
		nilErrorExpected(t, err)
	}
}
