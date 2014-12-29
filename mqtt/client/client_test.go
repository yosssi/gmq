package client

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

var errTest = errors.New("test error")

type packetErr struct{}

func (p *packetErr) WriteTo(w io.Writer) (int64, error) {
	return 0, errTest
}

func (p *packetErr) Type() (byte, error) {
	return 0x00, errTest
}

func TestClient_Connect_ErrAlreadyConnected(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	if err := cli.Connect(nil); err != ErrAlreadyConnected {
		invalidError(t, err, ErrAlreadyConnected)
	}
}

func TestClient_Connect_newConnectionErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.Connect(nil); err == nil {
		notNilErrorExpected(t)
	}
}

func TestClient_Connect_restoreSession_TypeErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.sess = newSession(false, []byte("cliendID"))

	cli.sess.sendingPackets[1] = &packetErr{}

	err := cli.Connect(&ConnectOptions{
		Network: "tcp",
		Address: testAddress,
	})
	if err != errTest {
		invalidError(t, err, errTest)
	}
}

func TestClient_Connect_restoreSession(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.sess = newSession(false, []byte("cliendID"))

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.sess.sendingPackets[1] = publish

	pubrel, err := packet.NewPUBREL(&packet.PUBRELOptions{
		PacketID: 2,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.sess.sendingPackets[2] = pubrel

	cli.sess.sendingPackets[3] = packet.NewPINGREQ()

	err = cli.Connect(&ConnectOptions{
		Network: "tcp",
		Address: testAddress,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_Connect_sendCONNECTErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network: "tcp",
		Address: testAddress,
	})
	if err != packet.ErrInvalidClientIDCleanSession {
		invalidError(t, err, packet.ErrInvalidClientIDCleanSession)
	}
}

func TestClient_Connect_CloseErr(t *testing.T) {
	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err = cli.Connect(&ConnectOptions{
		Network: "tcp",
		Address: "localhost:1883",
	})

	if err != packet.ErrInvalidClientIDCleanSession {
		invalidError(t, err, packet.ErrInvalidClientIDCleanSession)
	}

	ln.Close()
}

func TestClient_Disconnect_sendEndDefault(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.conn.sendEnd <- struct{}{}

	time.Sleep(3 * time.Second)

	cli.conn.sendEnd <- struct{}{}

	if err := cli.Disconnect(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_Publish_connNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.Publish(nil); err != ErrNotYetConnected {
		invalidError(t, err, ErrNotYetConnected)
	}
}

func TestClient_Publish_newPUBLISHPacketErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("cliendID"))

	err := cli.Publish(&PublishOptions{
		QoS: 0x03,
	})

	if err != packet.ErrInvalidQoS {
		invalidError(t, err, packet.ErrInvalidQoS)
	}
}

func TestClient_Publish(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.conn.send = make(chan packet.Packet, 1)

	if err := cli.Publish(nil); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_Subscribe_connNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.Subscribe(nil); err != ErrNotYetConnected {
		invalidError(t, err, ErrNotYetConnected)
	}
}

func TestClient_Subscribe_ErrInvalidNoSubReq(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	if err := cli.Subscribe(nil); err != packet.ErrInvalidNoSubReq {
		invalidError(t, err, packet.ErrInvalidNoSubReq)
	}
}

func TestClient_Subscribe_generatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	id := minPacketID

	for {
		cli.sess.sendingPackets[id] = nil

		if id == maxPacketID {
			break
		}

		id++
	}

	err := cli.Subscribe(&SubscribeOptions{
		SubReqs: []*SubReq{
			&SubReq{},
		},
	})

	if err != ErrPacketIDExhaused {
		invalidError(t, err, ErrPacketIDExhaused)
	}
}

func TestClient_Subscribe_NewSUBSCRIBEErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	err := cli.Subscribe(&SubscribeOptions{
		SubReqs: []*SubReq{
			&SubReq{},
		},
	})

	if err != packet.ErrNoTopicFilter {
		invalidError(t, err, packet.ErrNoTopicFilter)
	}
}

func TestClient_Subscribe(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	cli.conn.unackSubs = make(map[string]MessageHandler)

	cli.conn.send = make(chan packet.Packet, 1)

	err := cli.Subscribe(&SubscribeOptions{
		SubReqs: []*SubReq{
			&SubReq{
				TopicFilter: []byte("topicFilter"),
			},
		},
	})

	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_Unsubscribe_connNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.Unsubscribe(nil); err != ErrNotYetConnected {
		invalidError(t, err, ErrNotYetConnected)
	}
}

func TestClient_Unsubscribe_ErrNoTopicFilter(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	if err := cli.Unsubscribe(nil); err != packet.ErrNoTopicFilter {
		invalidError(t, err, packet.ErrNoTopicFilter)
	}
}

func TestClient_Unsubscribe_generatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	id := minPacketID

	for {
		cli.sess.sendingPackets[id] = nil

		if id == maxPacketID {
			break
		}

		id++
	}

	err := cli.Unsubscribe(&UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte{0x00},
		},
	})

	if err != ErrPacketIDExhaused {
		invalidError(t, err, ErrPacketIDExhaused)
	}
}

func TestClient_Unsubscribe_NewUNSUBSCRIBEErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	err := cli.Unsubscribe(&UnsubscribeOptions{
		TopicFilters: [][]byte{
			make([]byte, 0),
		},
	})

	if err != packet.ErrNoTopicFilter {
		invalidError(t, err, packet.ErrNoTopicFilter)
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.sess = newSession(false, []byte("clientID"))

	cli.conn.send = make(chan packet.Packet, 1)

	err := cli.Unsubscribe(&UnsubscribeOptions{
		TopicFilters: [][]byte{
			[]byte{0x00},
		},
	})

	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_Terminate(t *testing.T) {
	cli := &Client{}

	cli.disconnEndc = make(chan struct{}, 1)

	cli.Terminate()
}

func TestClient_send_connNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.send(nil); err != ErrNotYetConnected {
		invalidError(t, err, ErrNotYetConnected)
	}
}

func TestClient_sendCONNECT_optsNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if err := cli.sendCONNECT(nil); err != packet.ErrInvalidClientIDCleanSession {
		invalidError(t, err, packet.ErrInvalidClientIDCleanSession)
	}
}

func TestClient_receive_connNil(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	if _, err := cli.receive(); err != ErrNotYetConnected {
		invalidError(t, err, ErrNotYetConnected)
	}
}

func TestClient_receive_ReadByteErr(t *testing.T) {
	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	c := make(chan struct{})

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			<-c
			conn.Write([]byte{packet.TypePUBACK << 4})
			conn.Close()
		}
	}()

	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err = cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	defer cli.Disconnect()

	c <- struct{}{}

	cli.conn.wg.Wait()

	if err := ln.Close(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_receive_ReadFullErr(t *testing.T) {
	ln, err := net.Listen("tcp", ":1884")
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	c := make(chan struct{})

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			<-c
			conn.Write([]byte{packet.TypePUBACK << 4, 0x80, 0x01})
			conn.Close()
		}
	}()

	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err = cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1884",
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	defer cli.Disconnect()

	c <- struct{}{}

	cli.conn.wg.Wait()

	if err := ln.Close(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_clean_cleanSess(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.sess = newSession(true, []byte("clientID"))

	cli.clean()
}

func TestClient_waitPacket_timeout(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	cli.conn.wg.Add(1)

	cli.waitPacket(nil, 1, errTest)
}

func TestClient_receivePackets_handlePacketErr(t *testing.T) {
	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	c := make(chan struct{})

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			<-c
			if _, err := conn.Write([]byte{packet.TypePUBACK << 4, 0x02, 0x00, 0x01}); err != nil {
				return
			}
		}
	}()

	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err = cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	defer cli.Disconnect()

	c <- struct{}{}

	cli.conn.wg.Wait()

	if err := ln.Close(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePacket_TypeErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	p := &packetErr{}

	if err := cli.handlePacket(p); err != errTest {
		invalidError(t, err, errTest)
	}
}

func TestClient_handlePacket_PUBLISH(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.handlePacket(p)
}

func TestClient_handlePacket_PUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBACK(&packet.PUBACKOptions{
		PacketID: 1,
	})

	cli.handlePacket(p)
}

func TestClient_handlePacket_PUBREC(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBREC(&packet.PUBRECOptions{
		PacketID: 1,
	})

	cli.handlePacket(p)
}

func TestClient_handlePacket_PUBREL(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBREL(&packet.PUBRELOptions{
		PacketID: 1,
	})

	cli.handlePacket(p)
}

func TestClient_handlePacket_PUBCOMP(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBCOMP(&packet.PUBCOMPOptions{
		PacketID: 1,
	})

	cli.handlePacket(p)
}

func TestClient_handlePacket_SUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewSUBACKFromBytes([]byte{packet.TypeSUBACK << 4, 0x03}, []byte{0x00, 0x01, 0x00})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.handlePacket(p)
}

func TestClient_handlePacket_UNSUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewUNSUBACKFromBytes([]byte{packet.TypeUNSUBACK << 4, 0x02}, []byte{0x00, 0x01})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.handlePacket(p)
}

func TestClient_handlePacket_ErrInvalidPacketType(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewCONNECT(&packet.CONNECTOptions{
		ClientID: []byte("cliendID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.handlePacket(p)
}

func TestClient_handleCONNACK_default(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.conn = &connection{}

	cli.conn.connack = make(chan struct{})

	cli.handleCONNACK()
}

func TestClient_handlePUBLISH_QoS0(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if err = cli.handlePUBLISH(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBLISH_QoS1_Err(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:      mqtt.QoS1,
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	p.(*packet.PUBLISH).PacketID = 0

	if err = cli.handlePUBLISH(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBLISH_QoS1(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:      mqtt.QoS1,
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if err = cli.handlePUBLISH(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBLISH_QoS2_ErrInvalidPacketID(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:      mqtt.QoS2,
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.sess.receivingPackets[1] = p

	if err = cli.handlePUBLISH(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBLISH_QoS2_NewPUBRECErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:      mqtt.QoS2,
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	p.(*packet.PUBLISH).PacketID = 0

	if err = cli.handlePUBLISH(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBLISH_QoS2(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:      mqtt.QoS2,
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	if err = cli.handlePUBLISH(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBACK_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBACK{
		PacketID: 1,
	}

	if err := cli.handlePUBACK(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBACK{
		PacketID: 1,
	}

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.sess.sendingPackets[1] = publish

	if err := cli.handlePUBACK(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBREC_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREC{
		PacketID: 1,
	}

	if err := cli.handlePUBREC(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBREC_NewPUBRELErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREC{
		PacketID: 0,
	}

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	publish.(*packet.PUBLISH).PacketID = 0

	cli.sess.sendingPackets[0] = publish

	if err := cli.handlePUBREC(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBREC(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREC{
		PacketID: 1,
	}

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.sess.sendingPackets[1] = publish

	if err := cli.handlePUBREC(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBREL_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREL{
		PacketID: 1,
	}

	if err := cli.handlePUBREL(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBREL_NewPUBCOMPErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREL{
		PacketID: 0,
	}

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	publish.(*packet.PUBLISH).PacketID = 0

	cli.sess.receivingPackets[0] = publish

	if err := cli.handlePUBREL(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBREL(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBREL{
		PacketID: 1,
	}

	publish, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.sess.receivingPackets[1] = publish

	if err := cli.handlePUBREL(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePUBCOMP_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBCOMP{
		PacketID: 1,
	}

	if err := cli.handlePUBCOMP(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handlePUBCOMP(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.PUBCOMP{
		PacketID: 1,
	}

	pubrel, err := packet.NewPUBREL(&packet.PUBRELOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
		return
	}

	cli.sess.sendingPackets[1] = pubrel

	if err := cli.handlePUBCOMP(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleSUBACK_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.SUBACK{
		PacketID:    1,
		ReturnCodes: []byte{mqtt.QoS0},
	}

	if err := cli.handleSUBACK(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handleSUBACK_ErrInvalidSUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.SUBACK{
		PacketID:    1,
		ReturnCodes: []byte{mqtt.QoS0},
	}

	subscribe, err := packet.NewSUBSCRIBE(&packet.SUBSCRIBEOptions{
		PacketID: 1,
		SubReqs: []*packet.SubReq{
			&packet.SubReq{
				TopicFilter: []byte("test"),
				QoS:         mqtt.QoS0,
			},
			&packet.SubReq{
				TopicFilter: []byte("test2"),
				QoS:         mqtt.QoS0,
			},
		},
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.sess.sendingPackets[1] = subscribe

	if err := cli.handleSUBACK(p); err != ErrInvalidSUBACK {
		invalidError(t, err, ErrInvalidSUBACK)
	}
}

func TestClient_handleSUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.SUBACK{
		PacketID:    1,
		ReturnCodes: []byte{packet.SUBACKRetFailure, mqtt.QoS0},
	}

	subscribe, err := packet.NewSUBSCRIBE(&packet.SUBSCRIBEOptions{
		PacketID: 1,
		SubReqs: []*packet.SubReq{
			&packet.SubReq{
				TopicFilter: []byte("test"),
				QoS:         mqtt.QoS0,
			},
			&packet.SubReq{
				TopicFilter: []byte("test2"),
				QoS:         mqtt.QoS0,
			},
		},
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.sess.sendingPackets[1] = subscribe

	if err := cli.handleSUBACK(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleUNSUBACK_validatePacketIDErr(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.UNSUBACK{
		PacketID: 1,
	}

	if err := cli.handleUNSUBACK(p); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_handleUNSUBACK(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	p := &packet.UNSUBACK{
		PacketID: 1,
	}

	unsubscribe, err := packet.NewUNSUBSCRIBE(&packet.UNSUBSCRIBEOptions{
		PacketID:     1,
		TopicFilters: [][]byte{[]byte("test")},
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.sess.sendingPackets[1] = unsubscribe

	if err := cli.handleUNSUBACK(p); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePINGRESP_ErrInvalidPINGRESP(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	if err := cli.handlePINGRESP(); err != ErrInvalidPINGRESP {
		invalidError(t, err, ErrInvalidPINGRESP)
	}
}

func TestClient_handlePINGRESP_default(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}))

	if err := cli.handlePINGRESP(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handlePINGRESP(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	defer cli.Disconnect()

	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}, 1))

	if err := cli.handlePINGRESP(); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleErrorAndDisconn_connNil(t *testing.T) {
	cli := New(nil)

	cli.handleErrorAndDisconn(errTest)
}

func TestClient_handleErrorAndDisconn_default(t *testing.T) {
	cli := &Client{
		errorHandler: func(_ error) {},
	}

	cli.conn = &connection{}

	cli.disconnc = make(chan struct{})

	cli.handleErrorAndDisconn(errTest)
}

func TestClient_handleErrorAndDisconn(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	cli.handleErrorAndDisconn(errTest)

	time.Sleep(1 * time.Second)
}

func TestClient_sendPackets_sendErr(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:  "tcp",
		Address:  testAddress,
		ClientID: []byte("clientID"),
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.conn.Close(); err != nil {
		nilErrorExpected(t, err)
	}

	cli.conn.send <- &packetErr{}

	cli.conn.wg.Wait()
}

func TestClient_sendPackets_keepAliveErr(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:   "tcp",
		Address:   testAddress,
		ClientID:  []byte("clientID"),
		KeepAlive: 2,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.conn.Close(); err != nil {
		nilErrorExpected(t, err)
	}

	cli.conn.wg.Wait()
}

func TestClient_sendPackets_sendEnd(t *testing.T) {
	cli := New(nil)

	err := cli.Connect(&ConnectOptions{
		Network:   "tcp",
		Address:   testAddress,
		ClientID:  []byte("clientID"),
		KeepAlive: 3,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	time.Sleep(5 * time.Second)

	cli.conn.muPINGRESPs.Lock()
	cli.conn.pingresps = append(cli.conn.pingresps, make(chan struct{}))
	cli.conn.muPINGRESPs.Unlock()

	cli.conn.sendEnd <- struct{}{}

	cli.muConn.Lock()
	wg := &cli.conn.wg
	cli.muConn.Unlock()

	wg.Wait()
}

func TestClient_newPUBLISHPacket_generatePacketID(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	id := minPacketID

	for {
		cli.sess.sendingPackets[id] = nil

		if id == maxPacketID {
			break
		}

		id++
	}

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: mqtt.QoS1,
	})

	if err != ErrPacketIDExhaused {
		invalidError(t, err, ErrPacketIDExhaused)
	}
}

func TestClient_newPUBLISHPacket_ErrInvalidQoS(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: 0x03,
	})

	if err != packet.ErrInvalidQoS {
		invalidError(t, err, packet.ErrInvalidQoS)
	}
}

func TestClient_newPUBLISHPacket(t *testing.T) {
	cli := New(nil)

	cli.sess = newSession(false, []byte("clientID"))

	_, err := cli.newPUBLISHPacket(&PublishOptions{
		QoS: mqtt.QoS1,
	})

	if err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_validatePacketID_notExist(t *testing.T) {
	cli := New(nil)

	if err := cli.validatePacketID(map[uint16]packet.Packet{}, 1, packet.TypePUBLISH); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_validatePacketID_ptypeErr(t *testing.T) {
	cli := New(nil)

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: &packetErr{}}, 1, packet.TypePUBLISH); err != errTest {
		invalidError(t, err, errTest)
	}
}

func TestClient_validatePacketID_ErrInvalidPacketID(t *testing.T) {
	cli := New(nil)

	p, err := packet.NewPUBACK(&packet.PUBACKOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: p}, 1, packet.TypePUBLISH); err != packet.ErrInvalidPacketID {
		invalidError(t, err, packet.ErrInvalidPacketID)
	}
}

func TestClient_validatePacketID(t *testing.T) {
	cli := New(nil)

	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		PacketID: 1,
	})
	if err != nil {
		nilErrorExpected(t, err)
	}

	if err := cli.validatePacketID(map[uint16]packet.Packet{1: p}, 1, packet.TypePUBLISH); err != nil {
		nilErrorExpected(t, err)
	}
}

func TestClient_handleMessage_continue(t *testing.T) {
	cli := New(nil)

	cli.conn = &connection{}

	cli.conn.ackedSubs = map[string]MessageHandler{
		"test": nil,
	}

	cli.handleMessage([]byte("test"), nil)
}

func TestClient_handleMessage(t *testing.T) {
	cli := New(nil)

	cli.conn = &connection{}

	cli.conn.ackedSubs = map[string]MessageHandler{
		"test": func(_, _ []byte) {},
	}

	cli.handleMessage([]byte("test"), nil)
}

func TestNew_optsNil(t *testing.T) {
	cli := New(nil)

	cli.disconnc <- struct{}{}

	time.Sleep(500 * time.Millisecond)

	cli.disconnEndc <- struct{}{}

	cli.wg.Wait()

}

func TestNew(t *testing.T) {
	cli := New(&Options{
		ErrorHandler: func(_ error) {},
	})

	cli.disconnc <- struct{}{}

	time.Sleep(500 * time.Millisecond)

	cli.disconnEndc <- struct{}{}

	cli.wg.Wait()
}

func Test_match(t *testing.T) {
	testCases := []struct {
		in struct {
			topicName   string
			topicFilter string
		}
		out bool
	}{
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/score/wimbledon",
				topicFilter: "sport/tennis/player1/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/player1/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player2",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/player1/ranking",
				topicFilter: "sport/tennis/+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis/",
				topicFilter: "sport/tennis/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/tennis",
				topicFilter: "sport/tennis/+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test",
				topicFilter: "+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/test",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis/test/test",
				topicFilter: "+/tennis/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "test/tennis2/",
				topicFilter: "+/tennis/#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport//player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/test/player1",
				topicFilter: "sport/+/player1",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "sport/player1",
				topicFilter: "sport/+/player1",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "/finance",
				topicFilter: "+",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "#",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "+/monitor/Clients",
			},
			out: false,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/test/test",
				topicFilter: "$SYS/#",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients",
				topicFilter: "$SYS/monitor/+",
			},
			out: true,
		},
		{
			in: struct {
				topicName   string
				topicFilter string
			}{
				topicName:   "$SYS/monitor/Clients/test",
				topicFilter: "$SYS/monitor/+",
			},
			out: false,
		},
	}

	for _, tc := range testCases {
		if got := match(tc.in.topicName, tc.in.topicFilter); got != tc.out {
			t.Errorf("got => %t, want => %t", got, tc.out)
		}
	}
}

func invalidError(t *testing.T, err, want error) {
	if err == nil {
		t.Errorf("err => nil, want => %q", want)
	} else {
		t.Errorf("err => %q, want => %q", err, want)
	}
}
