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
