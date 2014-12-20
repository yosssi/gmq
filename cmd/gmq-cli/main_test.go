package main

import (
	"io"
	"strings"
	"testing"
)

func TestMain_printVersion(t *testing.T) {
	defer func(exitBak func(int)) {
		exit = exitBak
	}(exit)

	exit = func(_ int) {}

	b := true

	defer func(flagVBak *bool) {
		flagV = flagVBak
	}(flagV)

	flagV = &b

	main()
}

func TestMain_cmdNameEmpty(t *testing.T) {
	defer func(stdinBak io.Reader) {
		stdin = stdinBak
	}(stdin)

	stdin = strings.NewReader(" ")

	main()
}

func TestMain_newCommandErr(t *testing.T) {
	defer func(stdinBak io.Reader) {
		stdin = stdinBak
	}(stdin)

	stdin = strings.NewReader("notExitCmdName")

	main()
}

func TestMain_runErr(t *testing.T) {
	defer func(stdinBak io.Reader) {
		stdin = stdinBak
	}(stdin)

	stdin = strings.NewReader(cmdNameConn)

	main()
}

func TestMain(t *testing.T) {
	defer func(stdinBak io.Reader) {
		stdin = stdinBak
	}(stdin)

	stdin = strings.NewReader(cmdNameHelp)

	main()
}

func Test_printSuccess(t *testing.T) {
	printSuccess(cmdNameConn)
}

func Test_printError_errCmdArgsParse(t *testing.T) {
	printError(errCmdArgsParse, false)
}

func Test_cmdNameArgs(t *testing.T) {
	cmdNameArgs("test test1  test2")
}

func Test_disconn(t *testing.T) {
	ctx := newContext()

	ctx.wgMain.Add(1)

	go disconn(ctx)

	ctx.disconn <- struct{}{}

	ctx.disconnEnd <- struct{}{}

	ctx.wgMain.Wait()
}
