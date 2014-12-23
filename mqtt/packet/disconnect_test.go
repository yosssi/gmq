package packet

import "testing"

func TestNewDISCONNECT(t *testing.T) {
	p := NewDISCONNECT()

	ptype, err := p.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypeDISCONNECT {
		t.Errorf("ptype => %X, want => %X", ptype, TypeDISCONNECT)
	}
}
