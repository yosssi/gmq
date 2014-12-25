package packet

import "testing"

func TestNewFromBytes_ptypeErr(t *testing.T) {
	if _, err := NewFromBytes([]byte{}, nil); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func TestNewFromBytes_CONNACK(t *testing.T) {
	if _, err := NewFromBytes([]byte{TypeCONNACK << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewFromBytes_PUBACK(t *testing.T) {
	if _, err := NewFromBytes([]byte{TypePUBACK << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewFromBytes_PUBREC(t *testing.T) {
	if _, err := NewFromBytes([]byte{TypePUBREC << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewFromBytes_PUBCOMP(t *testing.T) {
	if _, err := NewFromBytes([]byte{TypePUBCOMP << 4, 0x02}, []byte{0x00, 0x00}); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewFromBytes_PINGRESP(t *testing.T) {
	if _, err := NewFromBytes([]byte{TypePINGRESP << 4, 0x00}, nil); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestNewFromBytes_ErrInvalidPacketType(t *testing.T) {
	if _, err := NewFromBytes([]byte{0x00 << 4}, nil); err != ErrInvalidPacketType {
		invalidError(t, err, ErrInvalidPacketType)
	}
}
