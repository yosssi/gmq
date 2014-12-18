package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// cmdHeader is the command-line input header.
const cmdHeader = "gmq-cli> "

// Global variables which are assigned other values while testing.
var (
	exit            = os.Exit
	stdin io.Reader = os.Stdin
)

// Command-line flag
var flagV = flag.Bool("v", false, "Print the version of GMQ Client and exit.")

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

	// Create a context.
	ctx := newContext()

	// Launch a goroutine which disconnects the Network Connection.
	go func() {
		for {
			select {
			case <-ctx.disconn:
				if err := disconnect(ctx); err != nil {
					printError(err)
				}
			}
		}
	}()

	// Create a scanner which reads lines from standard input.
	scanner := bufio.NewScanner(stdin)

	for printHeader(); scanner.Scan(); printHeader() {
		// Get the command name and the command arguments from
		// the scanner.
		cmdName, cmdArgs := cmdNameArgs(scanner.Text())

		// Skip the remaining processes if the command name is zero value.
		if cmdName == "" {
			continue
		}

		// Create a command.
		cmd, err := newCommand(cmdName, cmdArgs, ctx)
		if err != nil {
			printError(err)
			continue
		}

		// Run the command.
		if err := cmd.run(); err != nil {
			printError(err)
			continue
		}

		// Print the successfule message.
		printSuccess(cmdName)
	}
}

// printHeader prints the command-line input header to standard output.
func printHeader() {
	os.Stdout.WriteString(cmdHeader)
}

// printSuccess prints the successful message to standard output.
func printSuccess(cmdName string) {
	// Do nothing if the command name is the help command's name.
	if cmdName == cmdNameHelp {
		return
	}

	os.Stdout.WriteString("command was executed successfully.\n")
}

// printError prints the error to the standard error.
func printError(err error) {
	// Do nothing is the error is errCmdArgsParse.
	if err == errCmdArgsParse {
		return
	}

	fmt.Fprintln(os.Stderr, err)

	// Print the help of the GMQ Client commands if the error is errInvalidCmdName.
	if err == errInvalidCmdName {
		fmt.Println()
		printHelp()
	}
}

// cmdNameArgs extracts the command name and the command arguments from
// the parameter string.
func cmdNameArgs(s string) (string, []string) {
	// Split the string into the tokens.
	tokens := strings.Split(strings.TrimSpace(s), " ")

	// Get the command name from the tokens.
	cmdName := tokens[0]

	// Get the command arguments from the tokens.
	cmdArgs := make([]string, 0, len(tokens[1:]))
	for _, t := range tokens[1:] {
		// Skip the remaining processes if the token is zero value.
		if t == "" {
			continue
		}

		// Set the token to the command arguments.
		cmdArgs = append(cmdArgs, t)
	}

	return cmdName, cmdArgs
}
