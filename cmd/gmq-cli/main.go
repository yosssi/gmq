package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yosssi/gmq/mqtt/client"
)

// cmdHeader is the command-line input header.
const cmdHeader = "gmq-cli> "

// Global variables which are assigned other values while testing.
var (
	exit            = os.Exit
	stdin io.Reader = os.Stdin
)

// Command-line flags
var (
	flagV = flag.Bool("v", false, "Print the version of GMQ Client and exit.")
)

func init() {
	flag.Parse()
}

func main() {
	// Print the version of GMQ Client and exit if "v" flag is true.
	if *flagV {
		printVersion()
		exit(0)
		return
	}

	// Create an MQTT client.
	cli := client.New(nil)

	// Create a scanner which reads lines from standard input.
	scanner := bufio.NewScanner(stdin)

	for printHeader(); scanner.Scan(); printHeader() {
		// Get a string from the scanner.
		s := strings.TrimSpace(scanner.Text())

		// Skip the remaining processes if the string is zero value.
		if s == "" {
			continue
		}

		cmdName, cmdArgs := cmdNameArgs(s)

		// Create a command.
		cmd, err := newCommand(cmdName, cmdArgs, cli)
		if err != nil {
			printError(err)
			continue
		}

		// Run the command.
		if err := cmd.run(); err != nil {
			printError(err)
			continue
		}
	}
}

// printVersion prints the version of GMQ Client to standard output.
func printVersion() {
	fmt.Printf("GMQ Client %s\n", version)
}

// printHeader prints the command-line input header to standard output.
func printHeader() {
	os.Stdout.WriteString(cmdHeader)
}

// printError prints the error to standard error.
func printError(err error) {
	// Do nothing if the error is errCmdArgsParse.
	if err == errCmdArgsParse {
		return
	}

	fmt.Fprintf(os.Stderr, "%s.\n", err)

	// Print the help of the GMQ Client commands if the error is errInvalidCmdName.
	if err == errInvalidCmdName {
		fmt.Println()
		printHelp()
	}
}

// Extract the command name and the command arguments from the parameter string.
func cmdNameArgs(s string) (string, []string) {
	// Split the string into tokens.
	tokens := strings.Split(s, " ")

	// Get a command name from the tokens.
	cmdName := tokens[0]

	// Get command arguments from the tokens.
	var cmdArgs []string
	for _, t := range tokens[1:] {
		// Trim the token
		t = strings.TrimSpace(t)

		// Skip the remaining processes if the token is zero value.
		if t == "" {
			continue
		}

		// Set the token to the command arguments.
		cmdArgs = append(cmdArgs, t)
	}

	return cmdName, cmdArgs
}
