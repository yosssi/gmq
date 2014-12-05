package packet

import (
	"io/ioutil"
	"testing"
)

func Test_connect_WriteTo_errWriteFixedHeader(t *testing.T) {
	p, _ := NewCONNECT(nil)

	if _, err := p.WriteTo(&errWriter{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func Test_connect_WriteTo_errWriteVariableHeader(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).FixedHeader = nil

	if _, err := p.WriteTo(&errWriter{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func Test_connect_WriteTo_errWritePayload(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).FixedHeader = nil
	p.(*CONNECT).VariableHeader = nil

	if _, err := p.WriteTo(&errWriter{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func Test_connect_WriteTo(t *testing.T) {
	p, _ := NewCONNECT(nil)

	if _, err := p.WriteTo(ioutil.Discard); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_connect_setFixedHeader_0xFF000000(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).Payload = make([]byte, 10000000)

	p.(*CONNECT).setFixedHeader()
}

func Test_connect_setFixedHeader_0x00FF0000(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).Payload = make([]byte, 100000)

	p.(*CONNECT).setFixedHeader()
}

func Test_connect_setFixedHeader_0x0000FF00(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).Payload = make([]byte, 1000)

	p.(*CONNECT).setFixedHeader()
}

func Test_connect_setFixedHeader_0x000000FF(t *testing.T) {
	p, _ := NewCONNECT(nil)
	p.(*CONNECT).Payload = make([]byte, 1)

	p.(*CONNECT).setFixedHeader()
}

func Test_connect_setPayload(t *testing.T) {
	p, _ := NewCONNECT(&CONNECTOptions{
		WillTopic:   "willTopic",
		WillMessage: "willMessage",
		UserName:    "userName",
		Password:    "password",
	})

	p.(*CONNECT).setPayload()
}

func Test_connect_connectFlags(t *testing.T) {
	p, _ := NewCONNECT(&CONNECTOptions{
		WillTopic:   "willTopic",
		WillMessage: "willMessage",
		WillRetain:  true,
		UserName:    "userName",
		Password:    "password",
	})

	p.(*CONNECT).connectFlags()
}

func TestNewCONNECT(t *testing.T) {
	if p, _ := NewCONNECT(nil); p == nil {
		t.Error("p => nil, want => not nil")
	}
}
