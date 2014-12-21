package packet

import (
	"bytes"
	"errors"
	"testing"
)

type errWriter struct{}

func (w *errWriter) Write(p []byte) (int, error) {
	return 0, errTest
}

var errTest = errors.New("testError")

func TestBase_WriteTo_err(t *testing.T) {
	b := Base{}

	if _, err := b.WriteTo(&errWriter{}); err != errTest {
		if err == nil {
			t.Errorf("err => nil, want => %q", errTest)
		} else {
			t.Errorf("err => %q, want => %q", err, errTest)
		}
	}
}

func TestBase_WriteTo(t *testing.T) {
	fh, vh, p := "fixedHeader", "variableHeader", "payload"

	b := Base{
		FixedHeader:    []byte(fh),
		VariableHeader: []byte(vh),
		Payload:        []byte(p),
	}

	var bf bytes.Buffer

	if _, err := b.WriteTo(&bf); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if get, want := bf.String(), fh+vh+p; get != want {
		t.Errorf("b.String() => %q, want => %q", get, want)
	}
}

func TestBase_Type_err(t *testing.T) {
	b := Base{}

	if _, err := b.Type(); err != ErrInvalidFixedHeaderLen {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrInvalidFixedHeaderLen)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrInvalidFixedHeaderLen)
		}
	}
}

func TestBase_Type(t *testing.T) {
	var srcPtype byte = 0x10

	b := Base{
		FixedHeader: []byte{srcPtype, 0x00},
	}

	ptype, err := b.Type()
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
		return
	}

	if want := srcPtype >> 4; ptype != want {
		t.Errorf("ptype => %X, want => %X", ptype, want)
	}
}

func TestBase_appendRemainingLength(t *testing.T) {
	p := &CONNECT{}

	p.VariableHeader = []byte{0x00}
	p.Payload = []byte{0x00, 0x00}

	p.appendRemainingLength()

	if got, want := int(p.FixedHeader[0]), len(p.VariableHeader)+len(p.Payload); got != want {
		t.Errorf("got => %d, want => %d", got, want)
	}
}

func Test_appendRemainingLength(t *testing.T) {
	testCase := []struct {
		in  uint32
		out []byte
	}{
		{
			in:  0xFF000000,
			out: []byte{255, 0, 0, 0},
		},
		{
			in:  0x00FF0000,
			out: []byte{255, 0, 0},
		},
		{
			in:  0x0000FF00,
			out: []byte{255, 0},
		},
		{
			in:  0x000000FF,
			out: []byte{255},
		},
	}

	for _, tc := range testCase {
		b := appendRemainingLength([]byte{}, tc.in)

		if len(b) != len(tc.out) {
			t.Errorf("len(b) => %d, want => %d", len(b), len(tc.out))
			continue
		}

		for i, bt := range b {
			if bt != tc.out[i] {
				t.Errorf("bt => %X, want => %X", bt, tc.out[i])
				continue
			}
		}
	}
}
