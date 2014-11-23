package client

import "testing"

func TestClient_Conn_optsNill(t *testing.T) {
	if err, want := New().Conn(nil), "dial tcp 127.0.0.1:1883: connection refused"; err == nil || err.Error() != want {
		t.Error("err.Error() => %q, want => %q", err, want)
	}
}

func TestClient_Conn_errAlreadyConnected(t *testing.T) {
	opts := &ConnOpts{
		Host: "test.mosquitto.org",
	}

	cli := New()

	if err := cli.Conn(opts); err != nil {
		t.Error("err => %q, want => nil", err)
	}

	if err := cli.Conn(opts); err != ErrAlreadyConnected {
		if err == nil {
			t.Error("err => nil, want => %q", ErrAlreadyConnected)
		} else {
			t.Error("err => %q, want => %q", err, ErrAlreadyConnected)
		}
	}
}

func TestClient_Conn(t *testing.T) {
	opts := &ConnOpts{
		Host: "test.mosquitto.org",
	}

	if err := New().Conn(opts); err != nil {
		t.Error("err => %q, want => nil", err)
	}
}

func TestNew(t *testing.T) {
	if cli := New(); cli == nil {
		t.Error("cli should not be nil")
	}
}
