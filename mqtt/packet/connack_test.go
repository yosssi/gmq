package packet

import "testing"

func TestNewCONNACKFromBytes_errCONNACKInvalidFixedHeaderLen(t *testing.T) {
	if _, err := NewCONNACKFromBytes(nil, []byte{0x00, 0x00}); err != ErrCONNACKInvalidFixedHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidFixedHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidFixedHeaderLen)
		}
	}
}

func TestNewCONNACKFromBytes_errCONNACKInvalidVariableHeaderLen(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x02}, nil); err != ErrCONNACKInvalidVariableHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidVariableHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidVariableHeaderLen)
		}
	}
}

func TestNewCONNACKFromBytes_errCONNACKInvalidFixedHeaderFirstByte(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x21, 0x02}, []byte{0x00, 0x00}); err != ErrCONNACKInvalidFixedHeaderFirstByte {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidFixedHeaderFirstByte)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidFixedHeaderFirstByte)
		}
	}
}

func TestNewCONNACKFromBytes_errCONNACKInvalidRemainingLength(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x03}, []byte{0x00, 0x00}); err != ErrCONNACKInvalidRemainingLength {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidRemainingLength)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidRemainingLength)
		}
	}
}

func TestNewCONNACKFromBytes_errCONNACKInvalidConnectAcknowledgeFlags(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x02}, []byte{0x02, 0x00}); err != ErrCONNACKInvalidConnectAcknowledgeFlags {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidConnectAcknowledgeFlags)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidConnectAcknowledgeFlags)
		}
	}
}

func TestNewCONNACKFromBytes_errCONNACKInvalidConnectReturnCode(t *testing.T) {
	if _, err := NewCONNACKFromBytes([]byte{0x20, 0x02}, []byte{0x00, 0x06}); err != ErrCONNACKInvalidConnectReturnCode {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrCONNACKInvalidConnectReturnCode)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrCONNACKInvalidConnectReturnCode)
		}
	}
}

func TestNewCONNACKFromBytes(t *testing.T) {
	type in struct {
		fixedHeader    []byte
		variableHeader []byte
	}

	testCases := []struct {
		in  in
		out CONNACK
	}{
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x00, CONNACKAccepted},
			},
			out: CONNACK{
				SessionPresent:    false,
				ConnectReturnCode: CONNACKAccepted,
			},
		},
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x01, CONNACKUnacceptableProtocolVersion},
			},
			out: CONNACK{
				SessionPresent:    true,
				ConnectReturnCode: CONNACKUnacceptableProtocolVersion,
			},
		},
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x00, CONNACKIdentifierRejected},
			},
			out: CONNACK{
				SessionPresent:    false,
				ConnectReturnCode: CONNACKIdentifierRejected,
			},
		},
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x00, CONNACKServerUnavailable},
			},
			out: CONNACK{
				SessionPresent:    false,
				ConnectReturnCode: CONNACKServerUnavailable,
			},
		},
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x00, CONNACKBadUserNameOrPassword},
			},
			out: CONNACK{
				SessionPresent:    false,
				ConnectReturnCode: CONNACKBadUserNameOrPassword,
			},
		},
		{
			in: in{
				fixedHeader:    []byte{0x20, 0x02},
				variableHeader: []byte{0x00, CONNACKNotAuthorized},
			},
			out: CONNACK{
				SessionPresent:    false,
				ConnectReturnCode: CONNACKNotAuthorized,
			},
		},
	}

	for _, tc := range testCases {
		p, err := NewCONNACKFromBytes(tc.in.fixedHeader, tc.in.variableHeader)
		if err != nil {
			t.Errorf("err => %q, want => nil", err)
			continue
		}

		ca := p.(*CONNACK)

		if ca.SessionPresent != tc.out.SessionPresent {
			t.Errorf(" p.SessionPresent => %t, want => %t", ca.SessionPresent, tc.out.SessionPresent)
			continue
		}

		if ca.ConnectReturnCode != tc.out.ConnectReturnCode {
			t.Errorf(" ca.ConnectReturnCode => %X, want => %X", ca.ConnectReturnCode, tc.out.ConnectReturnCode)
		}
	}
}
