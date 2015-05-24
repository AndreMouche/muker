package server

import (
	"errors"
	"fmt"
	"github.com/openinx/muker/proto"
	"github.com/openinx/muker/utils"
	"io"
	"net"
	"time"
)

type Session struct {
	Conn       net.Conn
	ConnId     uint32
	SequenceId uint8
}

func NewSession(c net.Conn, connId uint32, sequenceId uint8) *Session {
	return &Session{
		Conn:       c,
		ConnId:     connId,
		SequenceId: sequenceId,
	}
}

func (s *Session) readPacket() (*proto.Packet, error) {
	header := make([]byte, 4)
	n, err := s.Conn.Read(header)

	if err != nil {
		return nil, err
	}

	if n != 4 {
		fmt.Printf("header length: %d\n", n)
		return nil, errors.New("Read packet header error : less than 4 bytes")
	}

	fmt.Printf("Read header: %x\n", header)

	pktLen := utils.BytesToUint24(header[:3])
	sequenceId := utils.BytesToUint8(header[3:])

	fmt.Printf("Read PacketLength: %d\n", pktLen)

	body := make([]byte, pktLen)
	n, err = s.Conn.Read(body)

	if err != nil && err != io.EOF {
		return nil, err
	}

	if n != pktLen {
		return nil, errors.New(fmt.Sprintf("Read packet body error: head length: %d, body length: %d", pktLen, n))
	}

	return proto.NewPacket(uint32(pktLen), sequenceId, body), nil
}
func (s *Session) HandleClient() {
	var pktBuf []byte
	var err error
	var written int

	fmt.Printf("RemoteAddr: %s\n", s.Conn.RemoteAddr().String())

	// send hande shake packet
	pkt := proto.DefaultHandShakePacket(s.ConnId)
	pktBuf, err = pkt.Write(s.SequenceId)
	if err != nil {
		defer s.Conn.Close()
		fmt.Printf("format packet to bytes error: %s\n", err.Error())
	}
	written, err = s.Conn.Write(pktBuf)
	fmt.Printf("send handshake pkt: %x\n", pktBuf)
	if err != nil || written != len(pktBuf) {
		defer s.Conn.Close()
		fmt.Printf("send hand shake pkt error: %s\n", err.Error())
	}

	// recv auth packet
	for {
		p, err := s.readPacket()
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		s.SequenceId = p.SequenceId + 1
		if authBuf, _ := proto.WritePacket(p.Buf, p.SequenceId); len(authBuf) > 0 {
			fmt.Printf("recv auth pkt: %x\n", authBuf)
		}
		pktBuf, err = proto.DefaultOkPacket().Write(s.SequenceId)
		s.Conn.Write(pktBuf)
		fmt.Printf("send auth ok pkt: %x\n", pktBuf)
		break
	}

	//recv command phase
	s.DoCommand()
}

func (s *Session) DoCommand() {

	defer s.Conn.Close()
	for {
		p, err := s.readPacket()

		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Printf("recv client pkt error: %s\n", err.Error())
			break
		}

		pBuf, _ := proto.WritePacket(p.Buf, p.SequenceId)

		// sequence increment
		s.SequenceId = p.SequenceId + 1
		fmt.Printf("recv client pkt : %x\n", pBuf)
		comType := p.Buf[0]
		supported, ok := proto.ComSupported[comType]

		// Command does not supported
		if !ok || !supported {
			pBuf, _ = proto.DefaultErrorPacket("Command Not Supported").Write(s.SequenceId)
			s.Conn.Write(pBuf)
			continue
		}

		switch comType {
		case proto.ComQuit:
			fmt.Printf("Command: Quit\n")
			s.Conn.Close()
		case proto.ComInitDB:
			fmt.Printf("Command: InitDB\n")
		case proto.ComQuery:
			s.doComQuery(p)
		case proto.ComFieldList:
			fmt.Printf("Command: FieldList\n")
		case proto.ComCreateDB:
			fmt.Printf("Command: CreateDB\n")
		case proto.ComDropDB:
			fmt.Printf("Command: DropDB\n")
		}
	}
}

func (s *Session) doComQuery(p *proto.Packet) {
	query := p.Buf[1:]
	fmt.Printf("Command Query: %s\n", query)
	pktBuf, _ := proto.DefaultOkPacket().Write(s.SequenceId)
	s.Conn.Write(pktBuf)
}