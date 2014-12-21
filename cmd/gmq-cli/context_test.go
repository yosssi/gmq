package main

import "testing"

func Test_context_initChan(t *testing.T) {
	(&context{}).initChan()
}

func Test_context_generatePacketID_errPacketIDExhaused(t *testing.T) {
	ctx := newContext()

	var i uint16

	for {
		ctx.packetIDs[i] = struct{}{}

		if i == maxPacketID {
			break
		}

		i++
	}

	if _, err := ctx.generatePacketID(); err != errPacketIDExhaused {
		errorfErr(t, err, errPacketIDExhaused)
	}
}

func Test_newContext(t *testing.T) {
	newContext()
}
