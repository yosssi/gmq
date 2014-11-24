package common

import "testing"

func TestNewNetworkConnection_DialErr(t *testing.T) {
	if _, err := NewNetworkConnection("", ""); err == nil {
		t.Error("err => nil, want => not nil")
	}
}

func TestNewNetworkConnection(t *testing.T) {
	if _, err := NewNetworkConnection("tcp", "test.mosquitto.org:1883"); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}
