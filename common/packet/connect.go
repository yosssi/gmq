package packet

import "io"

// Length of the Variable Header of the CONNECT Packet.
const lenCONNECTVariableHeader = 10

// connect represents a CONNECT Packet.
type connect struct {
	base
	// clientID is the Client Identifier (ClientId) identifies the Client to the Server.
	clientID string
	// cleanSession is the Clean Session of the connect flags.
	cleanSession bool
	// willTopic is the Will Topic of the payload.
	willTopic string
	// willMessage is the Will Message of the payload.
	willMessage string
	// willQoS is the Will QoS of the connect flags.
	willQoS uint
	// willRetain is the Will Retain of the connect flags.
	willRetain bool
	// userName is the user name used by the server for authentication and authorization.
	userName string
	// password is the password used by the server for authentication and authorization.
	password string
	// keepAlive is the Keep Alive in the variable header.
	keepAlive uint
}

// WriteTo writes the Packet data to the writer.
func (p *connect) WriteTo(w io.Writer) (int64, error) {
	var written int

	// Write the Fixed header.
	n, err := w.Write(p.fixedHeader)
	if err != nil {
		return int64(written), err
	}
	written += n

	// Write the Variable header.
	n, err = w.Write(p.variableHeader)
	if err != nil {
		return int64(written), err
	}
	written += n

	// Write the Payload.
	n, err = w.Write(p.payload)
	if err != nil {
		return int64(written), err
	}
	written += n

	return int64(written), nil
}

// setFixedHeader sets the Fixed header to the Packet.
func (p *connect) setFixedHeader() {
	// Create a byte slice holding the Fixed header.
	b := []byte{TypeCONNECT << 4}

	// Calculate the Remaining Length.
	rl := encodeLength(uint(lenCONNECTVariableHeader + len(p.payload)))

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

	p.fixedHeader = b
}

// setVariableHeader sets the Variable header.
func (p *connect) setVariableHeader() {
	// Create a byte slice holding the Variable header.
	b := make([]byte, lenCONNECTVariableHeader)

	// Set bytes.
	b[0] = 0                // Length MSB (0)
	b[1] = 4                // Length LSB (4)
	b[2] = 77               // 'M'
	b[3] = 81               // 'Q'
	b[4] = 84               // 'T'
	b[5] = 84               // 'T'
	b[6] = 4                // Level (4)
	b[7] = p.connectFlags() // Connect Flags

	// Set the Keep Alive.
	keepAlive := encodeUint16(uint16(p.keepAlive))
	b[8] = keepAlive[0]
	b[9] = keepAlive[1]

	// Set the byte slice to the Variable header.
	p.variableHeader = b
}

// setPayload sets the Payload to the Packet.
func (p *connect) setPayload() {
	// Create a byte slice holding the Payload.
	var b []byte

	// Append the Client Identifier.
	appendCONNECTPayload(&b, p.clientID)

	// Append the Will Topic and Will Message
	if p.will() {
		appendCONNECTPayload(&b, p.willTopic)
		appendCONNECTPayload(&b, p.willMessage)
	}

	if p.userName != "" {
		appendCONNECTPayload(&b, p.userName)
	}

	if p.password != "" {
		appendCONNECTPayload(&b, p.password)
	}

	p.payload = b
}

// connectFlags creates and returns the bytes representing the Connect Flags.
func (p *connect) connectFlags() byte {
	var b byte

	if p.userName != "" {
		b |= 128
	}

	if p.password != "" {
		b |= 64
	}

	if p.willRetain {
		b |= 32
	}

	b |= byte(p.willQoS) << 3

	if p.will() {
		b |= 4
	}

	if p.cleanSession {
		b |= 2
	}

	return b
}

// will returns if the Packet has both the Will Topic and Will Message.
func (p *connect) will() bool {
	return p.willTopic != "" && p.willMessage != ""
}

func appendCONNECTPayload(b *[]byte, s string) {
	*b = append(*b, encodeUint16(uint16(len(s)))...)
	*b = append(*b, []byte(s)...)
}

// NewCONNECT creates and returns a CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) Packet {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}
	opts.Init()

	// Create a CONNECT Packet.
	p := &connect{
		clientID:     opts.ClientID,
		cleanSession: *opts.CleanSession,
		willTopic:    opts.WillTopic,
		willMessage:  opts.WillMessage,
		willQoS:      opts.WillQoS,
		willRetain:   opts.WillRetain,
		userName:     opts.UserName,
		password:     opts.Password,
		keepAlive:    *opts.KeepAlive,
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
