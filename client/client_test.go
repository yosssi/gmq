package client

import "testing"

func TestClient_Connect_optsNill(t *testing.T) {
	if err := New().Connect(nil); err == nil {
		t.Errorf("err => nil, want => %q", err)
	}
}

func TestClient_Connect_errAlreadyConnected(t *testing.T) {
	opts := &ConnectOpts{
		Host: "test.mosquitto.org",
	}

	cli := New()

	if err := cli.Connect(opts); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	if err := cli.Connect(opts); err != ErrAlreadyConnected {
		if err == nil {
			t.Errorf("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Errorf("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Connect(t *testing.T) {
	opts := &ConnectOpts{
		Host: "test.mosquitto.org",
	}

	if err := New().Connect(opts); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func TestNew(t *testing.T) {
	if cli := New(); cli == nil {
		t.Error("cli should not be nil")
	}
}
