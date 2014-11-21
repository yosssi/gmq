package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yosssi/gmq/client"
	"github.com/yosssi/gmq/cmd/gmq-cli/cmd"
)

// header is a command line input header.
const header = "gmq-cli> "

// Special characters
const (
	space = " "
)

// Global variables which are injected with other values while testing
var (
	exit            = os.Exit
	stdin io.Reader = os.Stdin
)

// Command-line flags
var (
	v = flag.Bool("v", false, "Print the version and exit.")
)

// Commands
var cmds = []*cmd.Cmd{
	cmd.Conn,
}

func init() {
	flag.Parse()
}

func main() {
	// Print the version and exit.
	if *v {
		fmt.Printf("GMQ Client %s\n", client.Version)
		exit(0)
		return
	}

	// Create an MQTT client.
	cli := client.New()

	// Read lines from the standard input.
	scanner := bufio.NewScanner(stdin)

InputLoop:
	for writeHeader(); scanner.Scan(); writeHeader() {
		// Get the input data from the standard input.
		s := strings.TrimSpace(scanner.Text())

		if len(s) < 1 {
			continue
		}

		tokens := strings.Split(s, space)

		cmdName := tokens[0]
		cmdArgs := tokens[1:]

		// Print the help if the command name "help" is specified.
		if cmdName == "help" {
			for _, c := range cmds {
				fmt.Fprintf(os.Stdout, "%-11s %s\n", c.Name, c.Usage)
			}
			continue
		}

		for _, c := range cmds {
			if cmdName == c.Name {
				if c.Flag.Parse(cmdArgs) != nil {
					// An error message is printed by "c.Flag.Parse" and
					// error handling and printing are refrained here.
					continue InputLoop
				}

				c.Run(cli, c)

				fmt.Println(cmdName + " executed!")

				continue InputLoop
			}
		}

		fmt.Fprintf(os.Stderr, "command %q was not found.\n", cmdName)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func writeHeader() {
	os.Stdout.WriteString(header)
}
