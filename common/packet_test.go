package common

import "testing"

func TestPacket_NewPacketCONNECT_optsNil(t *testing.T) {
	NewPacketCONNECT(nil)
}

func TestPacket_NewPacketCONNECT(t *testing.T) {
	NewPacketCONNECT(&OptionsPacketCONNECT{})
}
