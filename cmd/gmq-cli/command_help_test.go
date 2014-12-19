package main

import "testing"

func Test_commandHelp_run(t *testing.T) {
	if err := (&commandHelp{}).run(); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_printHelp(t *testing.T) {
	printHelp()
}
