package packet

import "testing"

func TestNewPINGREQ(t *testing.T) {
	p := NewPINGREQ()

	ptype, err := p.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypePINGREQ {
		t.Errorf("ptype => %X, want => %X", ptype, TypePINGREQ)
	}
}
