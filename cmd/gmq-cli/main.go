package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yosssi/gmq/client"
)

// header is a command line input header.
const header = "gmq-cli> "

// Global variables which are injected with other values while testing
var (
	exit            = os.Exit
	stdin io.Reader = os.Stdin

	v = flag.Bool("v", false, "Print the version and exit.")
)

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

	scanner := bufio.NewScanner(stdin)

	os.Stdout.WriteString(header)

	for scanner.Scan() {
		fmt.Println(scanner.Text())

		os.Stdout.WriteString(header)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
