package packet

import "errors"

// Length of the Variable header of the CONNECT Packet.
const lenCONNECTVariableHeader = 10

// Protocol level
const protocolLevelVersion3_1_1 = 0x04

// Error values
var ErrCONNECTClientIDEmpty = errors.New("the Client Identifier is empty")

// CONNECT represents the CONNECT Packet.
type CONNECT struct {
	Base
	// ClientID is the Client Identifier which identifies the Client to the Server.
	ClientID string
	// CleanSession is the Clean Session of the Connect Flags.
	CleanSession bool
	// WillTopic is the Will Topic of the Payload.
	WillTopic string
	// WillMessage is the Will Message of the Payload.
	WillMessage string
	// WillQoS is the Will QoS of the Connect Flags.
	WillQoS uint
	// WillRetain is the Will Retain of the Connect Flags.
	WillRetain bool
	// UserName is the User Name used by the Server for authentication and authorization.
	UserName string
	// Password is the Password used by the Server for authentication and authorization.
	Password string
	// KeepAlive is the Keep Alive in the Variable header.
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

// setVariableHeader sets the Variable header to the Packet.
func (p *CONNECT) setVariableHeader() {
	// Create a byte slice holding the Variable header.
	b := make([]byte, lenCONNECTVariableHeader)

	// Set bytes.
	b[0] = 0x00                      // Length MSB (0)
	b[1] = 0x04                      // Length LSB (4)
	b[2] = 0x4D                      // 'M'
	b[3] = 0x51                      // 'Q'
	b[4] = 0x54                      // 'T'
	b[5] = 0x54                      // 'T'
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
	b = appendCONNECTPayload(b, p.ClientID)

	// Append the Will Topic and Will Message
	if p.will() {
		b = appendCONNECTPayload(b, p.WillTopic)
		b = appendCONNECTPayload(b, p.WillMessage)
	}

	if p.UserName != "" {
		b = appendCONNECTPayload(b, p.UserName)
	}

	if p.Password != "" {
		b = appendCONNECTPayload(b, p.Password)
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

// appendCONNECTPayload appends the length and the content of "s" to "b" and
// return the slice.
func appendCONNECTPayload(b []byte, s string) []byte {
	bytes := append(b, encodeUint16(uint16(len(s)))...)
	bytes = append(bytes, []byte(s)...)
	return bytes
}

// NewCONNECT creates and returns the CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}
	opts.Init()

	// Check the Client Identifier.
	if opts.ClientID == "" {
		return nil, ErrCONNECTClientIDEmpty
	}

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

	// Set the Variable header to the Packet.
	p.setVariableHeader()

	// Set the Payload to the Packet.
	p.setPayload()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}
