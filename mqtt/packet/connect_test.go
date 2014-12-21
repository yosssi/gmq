package packet

import "testing"

func TestCONNECT_setFixedHeader(t *testing.T) {
	p := &CONNECT{}

	p.setVariableHeader()

	p.setFixedHeader()

	fh := p.FixedHeader

	want := []byte{TypeCONNECT << 4, 0x0A}

	if len(fh) != len(want) {
		t.Errorf("len(fh) => %d, want => %d", len(fh), len(want))
		return
	}

	for i, b := range fh {
		if b != want[i] {
			t.Errorf("b => %X, want => %X", b, want[i])
			return
		}
	}
}

func TestCONNECT_setVariableHeader(t *testing.T) {
	p := &CONNECT{}

	p.setVariableHeader()

	vh := p.VariableHeader

	want := []byte{0, 4, 77, 81, 84, 84, 4, 0, 0, 0}

	if len(vh) != len(want) {
		t.Errorf("len(vh) => %d, want => %d", len(vh), len(want))
		return
	}

	for i, b := range vh {
		if b != want[i] {
			t.Errorf("b => %X, want => %X", b, want[i])
			return
		}
	}
}

func TestCONNECT_setPayload(t *testing.T) {
	p := &CONNECT{
		ClientID:    "clientID",
		WillTopic:   "willTopic",
		WillMessage: "willMessage",
		UserName:    "userName",
		Password:    "password",
	}

	p.setPayload()

	pl := p.Payload

	want := []byte{0, 8, 99, 108, 105, 101, 110, 116, 73, 68, 0, 9, 119, 105, 108, 108, 84, 111, 112, 105, 99, 0, 11, 119, 105, 108, 108, 77, 101, 115, 115, 97, 103, 101, 0, 8, 117, 115, 101, 114, 78, 97, 109, 101, 0, 8, 112, 97, 115, 115, 119, 111, 114, 100}

	if len(pl) != len(want) {
		t.Errorf("len(pl) => %d, want => %d", len(pl), len(want))
		return
	}

	for i, b := range pl {
		if b != want[i] {
			t.Errorf("b => %X, want => %X", b, want[i])
			return
		}
	}
}
func TestCONNECT_connectFlags(t *testing.T) {
	p := &CONNECT{
		ClientID:     "clientID",
		WillTopic:    "willTopic",
		WillMessage:  "willMessage",
		WillRetain:   true,
		UserName:     "userName",
		Password:     "password",
		CleanSession: true,
	}

	if got, want := p.connectFlags(), byte(0XE6); got != want {
		t.Errorf("p.connectFlags() => %X, want => %X", got, want)
	}
}

func TestCONNECT_will(t *testing.T) {
	testCase := []struct {
		in  *CONNECT
		out bool
	}{
		{
			in: &CONNECT{
				WillTopic:   "",
				WillMessage: "",
			},
			out: false,
		},
		{
			in: &CONNECT{
				WillTopic:   "",
				WillMessage: "willMessage",
			},
			out: false,
		},
		{
			in: &CONNECT{
				WillTopic:   "willTopic",
				WillMessage: "",
			},
			out: false,
		},
		{
			in: &CONNECT{
				WillTopic:   "willTopic",
				WillMessage: "willMessage",
			},
			out: true,
		},
	}

	for _, tc := range testCase {
		if got := tc.in.will(); got != tc.out {
			t.Errorf("tc.in.will() => %t, want => %t", got, tc.out)
			continue
		}
	}
}

func TestCONNECT_appendCONNECTPayload(t *testing.T) {
	b := appendCONNECTPayload([]byte{}, "test")

	want := []byte{0, 4, 116, 101, 115, 116}

	if len(b) != len(want) {
		t.Errorf("len(b) => %d, want => %d", len(b), len(want))
		return
	}

	for i, bt := range b {
		if bt != want[i] {
			t.Errorf("bt => %X, want => %X", bt, want[i])
			return
		}
	}
}

func TestNewCONNECT_optsNil(t *testing.T) {
	if _, err := NewCONNECT(nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNewCONNECT_errCONNECTInvalidWillQoS(t *testing.T) {
	if _, err := NewCONNECT(&CONNECTOptions{WillQoS: 3}); err != ErrCONNECTInvalidWillQoS {
		errorfErr(t, err, ErrCONNECTInvalidWillQoS)
	}
}

func TestNewCONNECT_errCONNECTWillTopicMessageEmpty(t *testing.T) {
	if _, err := NewCONNECT(&CONNECTOptions{WillTopic: "willTopic"}); err != ErrCONNECTWillTopicMessageEmpty {
		errorfErr(t, err, ErrCONNECTWillTopicMessageEmpty)
	}
}

func errorfErr(t *testing.T, err error, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)

	} else {
		t.Errorf("err => %q, want => %q", err, want)

	}
}
