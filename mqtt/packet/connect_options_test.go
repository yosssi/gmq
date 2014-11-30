package packet

import "testing"

func TestCONNECTOptions_Init(t *testing.T) {
	opts := &CONNECTOptions{}
	opts.Init()

	if *opts.CleanSession != DefaultCleanSession {
		t.Errorf("*opts.CleanSession => %t, want %t", *opts.CleanSession, DefaultCleanSession)
	}

	if *opts.KeepAlive != DefaultKeepAlive {
		t.Errorf("*opts.KeepAlive => %q, want %q", *opts.KeepAlive, DefaultKeepAlive)
	}
}
