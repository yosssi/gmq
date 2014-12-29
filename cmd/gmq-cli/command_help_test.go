package main

import "testing"

func Test_commandHelp_run(t *testing.T) {
	cmd := &commandHelp{}

	if err := cmd.run(); err != nil {
		nilErrorExpected(t, err)
	}
}

func Test_newCommandHelp(t *testing.T) {
	if cmd := newCommandHelp(); cmd == nil {
		t.Error("cmd => nil, want => not nil")
	}
}

func Test_printHelp(t *testing.T) {
	printHelp()
}

func nilErrorExpected(t *testing.T, err error) {
	t.Errorf("err => %q, want => nil", err)
}
