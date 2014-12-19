package main

import "testing"

func Test_context_initChan(t *testing.T) {
	(&context{}).initChan()
}

func Test_newContext(t *testing.T) {
	newContext()
}
