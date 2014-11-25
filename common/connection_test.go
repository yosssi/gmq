package common

import "testing"

const testAddress = "iot.eclipse.org:1883"

func TestNewConnection_DialErr(t *testing.T) {
	if _, err := NewConnection("", ""); err == nil {
		t.Error("err => nil, want => not nil")
	}
}

func TestNewConnection(t *testing.T) {
	if _, err := NewConnection("tcp", testAddress); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}
