package common

// MQTT Control Packet types
const (
	MQTTControlPacketTypeCONNECT     = 1
	MQTTControlPacketTypeCONNACK     = 2
	MQTTControlPacketTypePUBLISH     = 3
	MQTTControlPacketTypePUBACK      = 4
	MQTTControlPacketTypePUBREC      = 5
	MQTTControlPacketTypePUBREL      = 6
	MQTTControlPacketTypePUBCOMP     = 7
	MQTTControlPacketTypeSUBSCRIBE   = 8
	MQTTControlPacketTypeSUBACK      = 9
	MQTTControlPacketTypeUNSUBSCRIBE = 10
	MQTTControlPacketTypeUNSUBACK    = 11
	MQTTControlPacketTypePINGREQ     = 12
	MQTTControlPacketTypePINGRESP    = 13
	MQTTControlPacketTypeDISCONNECT  = 14
)
