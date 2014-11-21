package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

func init() {
	exit = func(_ int) {}
}

func TestMain_printVersion(t *testing.T) {
	defer func(b bool) {
		*v = b
	}(*v)

	*v = true

	main()
}
func TestMain(t *testing.T) {
	defer func(orig io.Reader) {
		stdin = orig
	}(stdin)

	stdin = strings.NewReader("test")

	main()
}

func TestMain_scannerErr(t *testing.T) {
	defer func(orig io.Reader) {
		stdin = orig
	}(stdin)

	stdin = bytes.NewReader(make([]byte, bufio.MaxScanTokenSize))

	main()
}
