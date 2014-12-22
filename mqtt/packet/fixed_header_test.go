package packet

import "testing"

func Test_fixedHeader_ptype_errInvalidFixedHeaderLen(t *testing.T) {
	var fh fixedHeader

	if _, err := fh.ptype(); err != ErrInvalidFixedHeaderLen {
		invalidError(t, err, ErrInvalidFixedHeaderLen)
	}
}

func Test_fixedHeader_ptype(t *testing.T) {
	fh := fixedHeader([]byte{TypeCONNECT << 4})

	ptype, err := fh.ptype()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypeCONNECT {
		t.Errorf("ptype => %X, want => %X", ptype, TypeCONNECT)
	}
}

func invalidError(t *testing.T, err, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
