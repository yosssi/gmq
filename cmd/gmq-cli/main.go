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

	scanner := bufio.NewScanner(stdin)

InputLoop:
	for printHeader(); scanner.Scan(); printHeader() {
		s := strings.TrimSpace(scanner.Text())

		if s == "" {
			continue
		}

		fmt.Println(s)

		continue InputLoop
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
