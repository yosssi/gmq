package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Length of the Variable header of the CONNECT Packet.
const lenCONNECTVariableHeader = 10

// Protocol level
const protocolLevelVersion3_1_1 = 0x04

// Error values
var (
	ErrCONNECTInvalidWillQoS        = errors.New("the Will QoS is invalid")
	ErrCONNECTWillTopicMessageEmpty = errors.New("the Will Topic or the Will Message is empty")
)

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

	// Append the Remaining Length to the slice and set it to the Fixed Header.
	p.FixedHeader = appendRemainingLength(b, rl)
}

// setVariableHeader sets the Variable header to the Packet.
func (p *CONNECT) setVariableHeader() {
	// Create a byte slice holding the Variable header.
	b := make([]byte, lenCONNECTVariableHeader)

	// Set bytes to the slice.
	b[0] = 0x00                      // Length MSB (0)
	b[1] = 0x04                      // Length LSB (4)
	b[2] = 0x4D                      // 'M'
	b[3] = 0x51                      // 'Q'
	b[4] = 0x54                      // 'T'
	b[5] = 0x54                      // 'T'
	b[6] = protocolLevelVersion3_1_1 // Protocol Level
	b[7] = p.connectFlags()          // Connect Flags

	// Set the Keep Alive to the slice.
	keepAlive := encodeUint16(uint16(p.KeepAlive))
	b[8] = keepAlive[0]
	b[9] = keepAlive[1]

	// Set the slice to the Variable header.
	p.VariableHeader = b
}

// setPayload sets the Payload to the Packet.
func (p *CONNECT) setPayload() {
	// Create a byte slice holding the Payload.
	var b []byte

	// Append the Client Identifier to the slice.
	b = appendCONNECTPayload(b, p.ClientID)

	// Append the Will Topic and Will Message to the slice.
	if p.will() {
		b = appendCONNECTPayload(b, p.WillTopic)
		b = appendCONNECTPayload(b, p.WillMessage)
	}

	// Append the User Name to the slice.
	if p.UserName != "" {
		b = appendCONNECTPayload(b, p.UserName)
	}

	// Append the Password to the slice.
	if p.Password != "" {
		b = appendCONNECTPayload(b, p.Password)
	}

	// Set the slice to the Payload.
	p.Payload = b
}

// connectFlags creates and returns the bytes representing the Connect Flags.
func (p *CONNECT) connectFlags() byte {
	// Define byte which represents the Connect Flags.
	var b byte

	// Set 1 to the Bit 7 if the Packets has the User Name.
	if p.UserName != "" {
		b |= 0x80
	}

	// Set 1 to the Bit 6 if the Packets has the Password.
	if p.Password != "" {
		b |= 0x40
	}

	// Set 1 to the Bit 5 if the Packets has the Will Retain.
	if p.WillRetain {
		b |= 0x20
	}

	// Set 00, 01 or 02 to the Bit 4 and 3 according to the QoS.
	b |= byte(p.WillQoS) << 3

	// Set 1 to the Bit 2 if the Packets has the Will Topic and the Will Message.
	if p.will() {
		b |= 0x04
	}
	// Set 1 to the Bit 1 if the Packets has the Clean Session.
	if p.CleanSession {
		b |= 0x02
	}

	// The Bit 0 is reserved and should be set to 0.

	// Return the byte which represents the Connect Flags.
	return b
}

// will returns true if the Packet has both the Will Topic and Will Message.
func (p *CONNECT) will() bool {
	return p.WillTopic != "" && p.WillMessage != ""
}

// appendCONNECTPayload appends the length and the content of the string to
// the slice and return the slice.
func appendCONNECTPayload(b []byte, s string) []byte {
	b = append(b, encodeUint16(uint16(len(s)))...)
	b = append(b, []byte(s)...)
	return b
}

// NewCONNECT creates and returns a CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}
	opts.Init()

	// Check the Will QoS.
	if !mqtt.ValidQoS(opts.WillQoS) {
		return nil, ErrCONNECTInvalidWillQoS
	}

	// Check the Will Topic and the Will Message.
	if (opts.WillTopic != "" && opts.WillMessage == "") || (opts.WillTopic == "" && opts.WillMessage != "") {
		return nil, ErrCONNECTWillTopicMessageEmpty
	}

	// Create a CONNECT Packet.
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
