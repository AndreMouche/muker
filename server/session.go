package server

import (
	"errors"
	"fmt"
	"github.com/AndreMouche/logging"
	"github.com/openinx/muker/pools"
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
	Backends   *pools.ConnPool
	SessionId  string
}

func NewSession(c net.Conn, connId uint32, sequenceId uint8, backends *pools.ConnPool) *Session {
	return &Session{
		Conn:       c,
		ConnId:     connId,
		SequenceId: sequenceId,
		Backends:   backends,
		SessionId:  fmt.Sprintf("session_%v_%v", time.Now().Unix(), connId),
	}
}

func (s *Session) readPacket() (*proto.Packet, error) {
	header := make([]byte, 4)
	n, err := s.Conn.Read(header)

	if err != nil {
		return nil, err
	}

	if n != 4 {
		logging.Error(s.SessionId, "header length:", n)
		return nil, errors.New("Read packet header error : less than 4 bytes")
	}

	logging.Debugf("%v == Read header: %x", s.SessionId, header)

	pktLen := utils.BytesToUint24(header[:3])
	sequenceId := utils.BytesToUint8(header[3:])

	logging.Debug(s.SessionId, "Read PacketLength:", pktLen)

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

	logging.Info(s.SessionId, "RemoteAddr:", s.Conn.RemoteAddr().String())

	// send hande shake packet
	pkt := proto.DefaultHandShakePacket(s.ConnId)
	pktBuf, err = pkt.Write(s.SequenceId)
	if err != nil {
		defer s.Conn.Close()
		logging.Error(s.SessionId, "format packet to bytes error:", err)
	}
	written, err = s.Conn.Write(pktBuf)
	logging.Debugf("%v send handshake pkt: %x", s.SessionId, pktBuf)
	if err != nil || written != len(pktBuf) {
		logging.Error(s.SessionId, "send hand shake pkt error:", err)
		s.Conn.Close()
		return
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
			logging.Debugf("%v recv auth pkt: %x", s.SessionId, authBuf)
		}
		pktBuf, err = proto.DefaultOkPacket().Write(s.SequenceId)
		s.Conn.Write(pktBuf)
		logging.Debugf("%v send auth ok pkt: %x", s.SessionId, pktBuf)
		break
	}

	//recv command phase
	s.DoCommand()
}

func (s *Session) DoCommand() {

	defer func() {
		logging.Info(s.SessionId, "DoCommand Finished")
	}()
	for {
		p, err := s.readPacket()

		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			logging.Error(s.SessionId, "recv client pkt error:", err)
			return
		}

		pBuf, _ := proto.WritePacket(p.Buf, p.SequenceId)

		// sequence increment
		s.SequenceId = p.SequenceId + 1
		logging.Debugf("%v recv client pkt : %x", s.SessionId, pBuf)
		comType := p.Buf[0]
		supported, ok := proto.ComSupported[comType]

		// Command does not supported
		if !ok || !supported {
			logging.Error(s.SessionId, "Command Type is Not Supported")
			pBuf, _ = proto.DefaultErrorPacket("Command Not Supported").Write(s.SequenceId)
			s.Conn.Write(pBuf)
			continue
		}

		switch comType {
		case proto.ComQuit:
			s.doComQuit(p)
			return
		case proto.ComInitDB:
			s.doComInitDB(p)
		case proto.ComQuery:
			s.doComQuery(p)
		case proto.ComFieldList:
			s.doComFieldList(p)
		case proto.ComCreateDB:
			s.doComCreateDB(p)
		case proto.ComDropDB:
			s.doComDropDB(p)
		}

	}
}

// To fix issue: https://github.com/openinx/muker/issues/1
func (s *Session) doComQuit(p *proto.Packet) {
	logging.Info(s.SessionId, "Command Quit")
	s.Conn.Close()
}

func (s *Session) doComInitDB(p *proto.Packet) {
	s.doInnerCommand("InitDB", p)
}

func (s *Session) doComQuery(p *proto.Packet) {
	s.doInnerCommand("Query", p)
}

func (s *Session) doComFieldList(p *proto.Packet) {
	s.doInnerCommand("ComFieldList", p)
}

func (s *Session) doComCreateDB(p *proto.Packet) {
	s.doInnerCommand("CreateDB", p)
}

func (s *Session) doComDropDB(p *proto.Packet) {
	s.doInnerCommand("DropDB", p)
}

func (s *Session) doInnerCommand(comName string, p *proto.Packet) {
	defer logging.Info(s.SessionId, "do command ", comName, " end")
	query := p.Buf[1:]
	logging.Debugf("%s Command %s: %s", s.SessionId, comName, query)

	c, err := s.Backends.Get()
	if err != nil {
		logging.Error(s.SessionId, "Get Conn Error:", err)
	}

	// reuse conn, put back to backend connection pool.
	defer func() {
		err = s.Backends.Put(c)
		if err != nil {
			logging.Error(s.SessionId, "Put back to conn pool failed:", err)
		}
	}()

	logging.Debug(s.SessionId, "Connect to backend Client Sucessful.")

	err2 := c.WriteCommandPacket(p, s.Conn)
	if err2 != nil {
		logging.Error(s.SessionId, "Error: ", err2)
	}

}
