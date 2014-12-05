package packet

import "testing"

func TestNewDISCONNECT(t *testing.T) {
	p := NewDISCONNECT()

	fh := p.(*DISCONNECT).FixedHeader

	if get, want := len(fh), 2; get != want {
		t.Errorf("len(fh) => %d, want => %d", get, want)
		return
	}

	if get, want := fh[0], byte(TypeDISCONNECT<<4); get != want {
		t.Errorf("fh[0] => %d, want => %d", get, want)
		return
	}

	if get, want := fh[1], byte(0x00); get != want {
		t.Errorf("fh[0] => %d, want => %d", get, want)
		return
	}
}
