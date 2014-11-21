package client

import "testing"

func TestNew(t *testing.T) {
	if cli := New(nil); cli == nil {
		t.Error("cli should not be nil")
	}
}
