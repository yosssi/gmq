package main

import "testing"

func Test_newCommand_errInvalidCmdName(t *testing.T) {
	if _, err := newCommand("invalidCmdName", nil, nil); err != errInvalidCmdName {
		errorfErr(t, err, errInvalidCmdName)
	}
}

func Test_newCommand_cmdNameConn(t *testing.T) {
	if _, err := newCommand(cmdNameConn, nil, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_newCommand_cmdNameDisconn(t *testing.T) {
	if _, err := newCommand(cmdNameDisconn, nil, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_newCommand_cmdNameHelp(t *testing.T) {
	if _, err := newCommand(cmdNameHelp, nil, nil); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func errorfErr(t *testing.T, err error, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
