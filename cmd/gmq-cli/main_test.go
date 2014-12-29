package main

import (
	"strings"
	"testing"
)

func Test_main_flagV(t *testing.T) {
	*flagV = true

	main()

	*flagV = false
}

func Test_main_cmdNameEmpty(t *testing.T) {
	stdin = strings.NewReader(" ")

	main()
}

func Test_main_newCommandErr(t *testing.T) {
	stdin = strings.NewReader("invalidCmdName")

	main()
}

func Test_main_runErr(t *testing.T) {
	stdin = strings.NewReader("conn")

	main()
}

func Test_main_run(t *testing.T) {
	stdin = strings.NewReader("help")

	main()
}

func Test_printSuccess(t *testing.T) {
	printSuccess(cmdNameHelp)

	printSuccess(cmdNameConn)
}

func Test_printError(t *testing.T) {
	printError(errCmdArgsParse)

	printError(errInvalidCmdName)
}

func Test_cmdNameArgs(t *testing.T) {
	cmdName, cmdArgs := cmdNameArgs("cmdName  cmdArgs")

	if cmdName != "cmdName" {
		t.Errorf("cmdName => %s, want => %s", cmdName, "cmdName")
	}

	want := []string{"cmdArgs"}

	if len(cmdArgs) != len(want) || cmdArgs[0] != want[0] {
		t.Errorf("cmdArgs => %v, want => %v", cmdArgs, want)
	}
}

func Test_errorHandler(t *testing.T) {
	errorHandler(errCmdArgsParse)
}
