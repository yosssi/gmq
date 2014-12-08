package client

import "testing"

func TestOptions_Init(t *testing.T) {
	opts := &Options{}

	opts.Init()

	if *opts.ConnTimeout != DefaultConnTimeout {
		t.Errorf("*opts.ConnTimeout => %d, want => %d", *opts.ConnTimeout, DefaultConnTimeout)
	}
}
