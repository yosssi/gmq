package packet

// Length of the fixed header of the PUBACK Packet
const lenPUBACKFixedHeader = 2

// Length of the variable header of the PUBACK Packet
const lenPUBACKVariableHeader = 2

// PUBACK represents a PUBACK Packet.
type PUBACK struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// NewPUBACKFromBytes creates a PUBACK Packet
// from the byte data and returns it.
func NewPUBACKFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePUBACKBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBACK Packet.
	p := &PUBACK{
		PacketID: packetID,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validatePUBACKBytes validates the fixed header and the variable header.
func validatePUBACKBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenPUBACKFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePUBACK {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenPUBACKVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenPUBACKVariableHeader {
		return ErrInvalidVariableHeaderLen
	}

	return nil
}
