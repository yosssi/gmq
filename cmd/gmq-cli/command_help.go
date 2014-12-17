package main

import "fmt"

// String format for the help
const strHelpFmt = "%-8s %s\n"

// commandHelp represents a help command.
type commandHelp struct{}

// run prints the help of the GMQ Client commands.
func (cmd *commandHelp) run() error {
	printHelp()

	return nil
}

// newCommandHelp creates and returns a help command.
func newCommandHelp() *commandHelp {
	return &commandHelp{}
}

func printHelp() {
	printVersion()
	fmt.Println("Usage:")
	fmt.Printf(strHelpFmt, cmdNameConn, "establish a Network Connection and send a CONNECT Packet to the Server")
	fmt.Printf(strHelpFmt, cmdNameDisconn, "send a DISCONNECT Packet to the Server and disconnect the Network Connection")
	fmt.Printf(strHelpFmt, cmdNameHelp, "print this help message")
}
