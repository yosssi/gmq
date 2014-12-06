package client

import "testing"

func TestConnectOptions_Init(t *testing.T) {
	opts := &ConnectOptions{}

	opts.Init()

	if opts.Network != DefaultNetwork {
		t.Errorf("opts.Network => %q, want => %q", opts.Network, DefaultNetwork)
		return
	}

	if opts.Address != DefaultAddress {
		t.Errorf("opts.Address => %q, want => %q", opts.Address, DefaultAddress)
		return
	}
}
