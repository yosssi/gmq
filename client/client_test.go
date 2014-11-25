package client

import "testing"

const testAddress = "iot.eclipse.org:1883"

func TestClient_Connect_addressEmpty(t *testing.T) {
	if err := New().Connect("", nil); err == nil {
		t.Errorf("err => nil, want => %q", err)
	}
}

func TestClient_Connect_errAlreadyConnected(t *testing.T) {
	cli := New()

	if err := cli.Connect(testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Connect(testAddress, nil); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Connect(t *testing.T) {
	if err := New().Connect(testAddress, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNew(t *testing.T) {
	if cli := New(); cli == nil {
		t.Error("cli should not be nil")
	}
}
