package client

import "testing"

func TestSessionOptions_Init(t *testing.T) {
	opts := &SessionOptions{}

	opts.Init()

	if *opts.CleanSession != DefaultCleanSession {
		t.Errorf("*opts.CleanSession => %t, want => %t", *opts.CleanSession, DefaultCleanSession)
		return
	}

	if opts.ClientID != hostname {
		t.Errorf("opts.ClientID => %q, want => %q", opts.ClientID, hostname)
		return
	}
}
