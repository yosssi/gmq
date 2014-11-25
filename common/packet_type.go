package common

// MQTT Control Packet types
const (
	PacketTypeCONNECT     = 1
	PacketTypeCONNACK     = 2
	PacketTypePUBLISH     = 3
	PacketTypePUBACK      = 4
	PacketTypePUBREC      = 5
	PacketTypePUBREL      = 6
	PacketTypePUBCOMP     = 7
	PacketTypeSUBSCRIBE   = 8
	PacketTypeSUBACK      = 9
	PacketTypeUNSUBSCRIBE = 10
	PacketTypeUNSUBACK    = 11
	PacketTypePINGREQ     = 12
	PacketTypePINGRESP    = 13
	PacketTypeDISCONNECT  = 14
)
