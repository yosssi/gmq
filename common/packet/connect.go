package packet

// Length of the Variable header of the CONNECT Packet.
const lenCONNECTVariableHeader = 10

// Protocol levels
const protocolLevelVersion3_1_1 = 0x04

// CONNECT represents the CONNECT Packet.
type CONNECT struct {
	Base
	// clientID is the Client Identifier (ClientId) identifies the Client to the Server.
	ClientID string
	// cleanSession is the Clean Session of the connect flags.
	CleanSession bool
	// willTopic is the Will Topic of the payload.
	WillTopic string
	// willMessage is the Will Message of the payload.
	WillMessage string
	// willQoS is the Will QoS of the connect flags.
	WillQoS uint
	// willRetain is the Will Retain of the connect flags.
	WillRetain bool
	// userName is the user name used by the server for authentication and authorization.
	UserName string
	// password is the password used by the server for authentication and authorization.
	Password string
	// keepAlive is the Keep Alive in the variable header.
	KeepAlive uint
}

// setFixedHeader sets the Fixed header to the Packet.
func (p *CONNECT) setFixedHeader() {
	// Create a byte slice holding the Fixed header.
	b := []byte{TypeCONNECT << 4}

	// Calculate the Remaining Length.
	rl := encodeLength(uint(lenCONNECTVariableHeader + len(p.Payload)))

	if rl&0xFF000000 > 0 {
		b = append(b, byte((rl&0xFF000000)>>24))
		b = append(b, byte((rl&0x00FF0000)>>16))
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	} else if rl&0x00FF0000 > 0 {
		b = append(b, byte((rl&0x00FF0000)>>16))
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	} else if rl&0x0000FF00 > 0 {
		b = append(b, byte((rl&0x0000FF00)>>8))
		b = append(b, byte(rl&0x000000FF))
	} else {
		b = append(b, byte(rl&0x000000FF))
	}

	p.FixedHeader = b
}

// setVariableHeader sets the Variable header.
func (p *CONNECT) setVariableHeader() {
	// Create a byte slice holding the Variable header.
	b := make([]byte, lenCONNECTVariableHeader)

	// Set bytes.
	b[0] = 0                         // Length MSB (0)
	b[1] = 4                         // Length LSB (4)
	b[2] = 77                        // 'M'
	b[3] = 81                        // 'Q'
	b[4] = 84                        // 'T'
	b[5] = 84                        // 'T'
	b[6] = protocolLevelVersion3_1_1 // Protocol Level
	b[7] = p.connectFlags()          // Connect Flags

	// Set the Keep Alive.
	keepAlive := encodeUint16(uint16(p.KeepAlive))
	b[8] = keepAlive[0]
	b[9] = keepAlive[1]

	// Set the byte slice to the Variable header.
	p.VariableHeader = b
}

// setPayload sets the Payload to the Packet.
func (p *CONNECT) setPayload() {
	// Create a byte slice holding the Payload.
	var b []byte

	// Append the Client Identifier.
	appendCONNECTPayload(&b, p.ClientID)

	// Append the Will Topic and Will Message
	if p.will() {
		appendCONNECTPayload(&b, p.WillTopic)
		appendCONNECTPayload(&b, p.WillMessage)
	}

	if p.UserName != "" {
		appendCONNECTPayload(&b, p.UserName)
	}

	if p.Password != "" {
		appendCONNECTPayload(&b, p.Password)
	}

	p.Payload = b
}

// connectFlags creates and returns the bytes representing the Connect Flags.
func (p *CONNECT) connectFlags() byte {
	var b byte

	if p.UserName != "" {
		b |= 128
	}

	if p.Password != "" {
		b |= 64
	}

	if p.WillRetain {
		b |= 32
	}

	b |= byte(p.WillQoS) << 3

	if p.will() {
		b |= 4
	}

	if p.CleanSession {
		b |= 2
	}

	return b
}

// will returns if the Packet has both the Will Topic and Will Message.
func (p *CONNECT) will() bool {
	return p.WillTopic != "" && p.WillMessage != ""
}

func appendCONNECTPayload(b *[]byte, s string) {
	*b = append(*b, encodeUint16(uint16(len(s)))...)
	*b = append(*b, []byte(s)...)
}

// NewCONNECT creates and returns the CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) Packet {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}
	opts.Init()

	// Create the CONNECT Packet.
	p := &CONNECT{
		ClientID:     opts.ClientID,
		CleanSession: *opts.CleanSession,
		WillTopic:    opts.WillTopic,
		WillMessage:  opts.WillMessage,
		WillQoS:      opts.WillQoS,
		WillRetain:   opts.WillRetain,
		UserName:     opts.UserName,
		Password:     opts.Password,
		KeepAlive:    *opts.KeepAlive,
	}

	// Set the Variable header.
	p.setVariableHeader()

	// Set the Payload.
	p.setPayload()

	// Set the Fixed header.
	p.setFixedHeader()

	// Return the CONNECT Packet.
	return p
}
