package common

// MQTT Control Packet types
const (
	MQTTControlPacketTypeCONNECT     MQTTControlPacketType = 1
	MQTTControlPacketTypeCONNACK     MQTTControlPacketType = 2
	MQTTControlPacketTypePUBLISH     MQTTControlPacketType = 3
	MQTTControlPacketTypePUBACK      MQTTControlPacketType = 4
	MQTTControlPacketTypePUBREC      MQTTControlPacketType = 5
	MQTTControlPacketTypePUBREL      MQTTControlPacketType = 6
	MQTTControlPacketTypePUBCOMP     MQTTControlPacketType = 7
	MQTTControlPacketTypeSUBSCRIBE   MQTTControlPacketType = 8
	MQTTControlPacketTypeSUBACK      MQTTControlPacketType = 9
	MQTTControlPacketTypeUNSUBSCRIBE MQTTControlPacketType = 10
	MQTTControlPacketTypeUNSUBACK    MQTTControlPacketType = 11
	MQTTControlPacketTypePINGREQ     MQTTControlPacketType = 12
	MQTTControlPacketTypePINGRESP    MQTTControlPacketType = 13
	MQTTControlPacketTypeDISCONNECT  MQTTControlPacketType = 14
)

// MQTTControlPacketType represents an MQTT Control Packet type
type MQTTControlPacketType byte
