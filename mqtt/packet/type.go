package packet

// MQTT Control Packet types
const (
	TypeCONNECT     = 0x01
	TypeCONNACK     = 0x02
	TypePUBLISH     = 0x03
	TypePUBACK      = 0x04
	TypePUBREC      = 0x05
	TypePUBREL      = 0x06
	TypePUBCOMP     = 0x07
	TypeSUBSCRIBE   = 0x08
	TypeSUBACK      = 0x09
	TypeUNSUBSCRIBE = 0x0A
	TypeUNSUBACK    = 0x0B
	TypePINGREQ     = 0x0C
	TypePINGRESP    = 0x0D
	TypeDISCONNECT  = 0x0E
)
