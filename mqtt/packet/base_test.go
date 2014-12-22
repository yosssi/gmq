package packet

import (
	"io/ioutil"
	"testing"
)

func Test_base_WriteTo(t *testing.T) {
	b := base{
		fixedHeader:    []byte{0x00},
		variableHeader: []byte{0x00, 0x00},
		payload:        []byte{0x00, 0x00, 0x00},
	}

	n, err := b.WriteTo(ioutil.Discard)
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if want := int64(len(b.fixedHeader) + len(b.variableHeader) + len(b.payload)); n != want {
		t.Errorf("n => %d, want => %d", n, want)
	}
}

func Test_base_Type(t *testing.T) {
	b := base{
		fixedHeader: []byte{TypeCONNECT << 4},
	}

	ptype, err := b.Type()
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if ptype != TypeCONNECT {
		t.Errorf("ptype => %X, want => %X", ptype, TypeCONNECT)
	}
}

func Test_base_appendRemainingLength(t *testing.T) {
	b := base{
		variableHeader: []byte{0x00},
		payload:        []byte{0x00, 0x00},
	}

	b.appendRemainingLength()

	if want := []byte{0x03}; len(b.fixedHeader) != len(want) || b.fixedHeader[0] != want[0] {
		t.Errorf("b.fixedHeader => %v, want => %v", b.fixedHeader, want)
	}
}

func Test_appendRemainingLength(t *testing.T) {
	var b []byte

	b = appendRemainingLength(b, 0xFF000000)

	want := []byte{0xFF, 0x00, 0x00, 0x00}

	if len(b) != len(want) || b[0] != want[0] || b[1] != want[1] || b[2] != want[2] || b[3] != want[3] {
		t.Errorf("b => %X, want => %X", b, want)
	}
}

func nilErrorExpected(t *testing.T, err error) {
	t.Errorf("err => %q, want => nil", err)
}
