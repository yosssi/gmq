package packet

import "testing"

func TestNewPINGREQ(t *testing.T) {
	fh := NewPINGREQ().(*PINGREQ).FixedHeader

	if get, want := len(fh), 2; get != want {
		t.Errorf("len(fh) => %d, want => %d", get, want)
		return
	}

	if get, want := fh[0], byte(TypePINGREQ<<4); get != want {
		t.Errorf("fh[0] => %d, want => %d", get, want)
		return
	}

	if get, want := fh[1], byte(0x00); get != want {
		t.Errorf("fh[0] => %d, want => %d", get, want)
		return
	}

}
