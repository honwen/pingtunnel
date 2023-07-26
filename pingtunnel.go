package pingtunnel

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"

	"github.com/esrrhs/pingtunnel/pkg/common"
	"github.com/esrrhs/pingtunnel/pkg/msg"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"google.golang.org/protobuf/proto"
)

func sendICMP(id int, sequence int, conn icmp.PacketConn, server *net.IPAddr, target string,
	connId string, msgType uint32, data []byte, sproto int, rproto int, key int,
	tcpmode int, tcpmode_buffer_size int, tcpmode_maxwin int, tcpmode_resend_time int, tcpmode_compress int, tcpmode_stat int,
	timeout int) {

	m := &msg.PingMsg{
		Id:                  connId,
		Type:                (int32)(msgType),
		Target:              target,
		Data:                data,
		Rproto:              (int32)(rproto),
		Key:                 (int32)(key),
		Tcpmode:             (int32)(tcpmode),
		TcpmodeBuffersize:   (int32)(tcpmode_buffer_size),
		TcpmodeMaxwin:       (int32)(tcpmode_maxwin),
		TcpmodeResendTimems: (int32)(tcpmode_resend_time),
		TcpmodeCompress:     (int32)(tcpmode_compress),
		TcpmodeStat:         (int32)(tcpmode_stat),
		Timeout:             (int32)(timeout),
		Magic:               (int32)(msg.PingMsg_MAGIC),
	}

	mb, err := proto.Marshal(m)
	if err != nil {
		log.Printf("sendICMP Marshal msg.PingMsg error %s %s", server.String(), err)
		return
	}

	body := &icmp.Echo{
		ID:   id,
		Seq:  sequence,
		Data: mb,
	}

	msg := &icmp.Message{
		Type: (ipv4.ICMPType)(sproto),
		Code: 0,
		Body: body,
	}

	bytes, err := msg.Marshal(nil)
	if err != nil {
		log.Printf("sendICMP Marshal error %s %s", server.String(), err)
		return
	}

	conn.WriteTo(bytes, server)
}

func recvICMP(workResultLock *sync.WaitGroup, exit *bool, conn icmp.PacketConn, recv chan<- *Packet) {

	defer common.CrashLog()

	(*workResultLock).Add(1)
	defer (*workResultLock).Done()

	bytes := make([]byte, 10240)
	for !*exit {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		n, srcaddr, err := conn.ReadFrom(bytes)

		if err != nil {
			nerr, ok := err.(net.Error)
			if !ok || !nerr.Timeout() {
				log.Printf("Error read icmp message %s", err)
				continue
			}
		}

		if n <= 0 {
			continue
		}

		echoId := int(binary.BigEndian.Uint16(bytes[4:6]))
		echoSeq := int(binary.BigEndian.Uint16(bytes[6:8]))

		my := &msg.PingMsg{}
		err = proto.Unmarshal(bytes[8:n], my)
		if err != nil {
			log.Printf("Unmarshal msg.PingMsg error: %s", err)
			continue
		}

		if my.Magic != (int32)(msg.PingMsg_MAGIC) {
			log.Printf("processPacket data invalid %s", my.Id)
			continue
		}

		recv <- &Packet{my: my,
			src:    srcaddr.(*net.IPAddr),
			echoId: echoId, echoSeq: echoSeq}
	}
}

type Packet struct {
	my      *msg.PingMsg
	src     *net.IPAddr
	echoId  int
	echoSeq int
}

const (
	FRAME_MAX_SIZE int = 888
	FRAME_MAX_ID   int = 1000000
)
