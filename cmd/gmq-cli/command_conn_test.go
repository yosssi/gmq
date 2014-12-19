package main

import "testing"

func Test_commandConn_run(t *testing.T) {
	ctx := newContext()

	cmd, err := newCommandConn(nil, ctx)
	if err != nil {
		t.Errorf("err => %q, want => nil", err)
	}

	cmd.run()
}
