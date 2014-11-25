package packet

import (
	"io/ioutil"
	"testing"
)

func Test_connect_WriteTo(t *testing.T) {
	p := NewCONNECT(nil)
	p.WriteTo(ioutil.Discard)
}

func TestNewCONNECT(t *testing.T) {
	if p := NewCONNECT(nil); p == nil {
		t.Error("p => nil, want => not nil")
	}
}
