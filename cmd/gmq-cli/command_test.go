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

func Test_newCommand_cmdNamePub(t *testing.T) {
	if _, err := newCommand(cmdNamePub, nil, newContext()); err != nil {
		t.Errorf("err => %q, want => nil", err)
	}
}

func Test_newCommand_cmdNamePub_errCmdArgsParse(t *testing.T) {
	if _, err := newCommand(cmdNamePub, []string{"-not-exist-flag"}, newContext()); err != errCmdArgsParse {
		errorfErr(t, err, errCmdArgsParse)
	}
}

func Test_newCommand_cmdNamePub_errPacketIDExhaused(t *testing.T) {
	ctx := newContext()

	var i uint16

	for {
		ctx.packetIDs[i] = struct{}{}

		if i == maxPacketID {
			break
		}

		i++
	}

	if _, err := newCommand(cmdNamePub, nil, ctx); err != errPacketIDExhaused {
		errorfErr(t, err, errPacketIDExhaused)
	}
}

func errorfErr(t *testing.T, err error, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
