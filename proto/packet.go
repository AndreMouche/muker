package proto

import (
	"github.com/openinx/muker/utils"
)

const (
	ComSleep            = 0x00
	ComQuit             = 0x01
	ComInitDB           = 0x02
	ComQuery            = 0x03
	ComFieldList        = 0x04
	ComCreateDB         = 0x05
	ComDropDB           = 0x06
	ComRefresh          = 0x07
	ComShutdown         = 0x08
	ComStatistics       = 0x09
	ComProcessInfo      = 0x0a
	ComConnect          = 0x0b
	ComProcessKill      = 0x0c
	ComDebug            = 0x0d
	ComPing             = 0x0e
	ComTime             = 0x0f
	ComDelayedInsert    = 0x10
	ComChangeUser       = 0x11
	ComBinlogDump       = 0x12
	ComTableDump        = 0x13
	ComConnectOut       = 0x14
	ComRegisterSlave    = 0x15
	ComStmtPrepare      = 0x16
	ComStmtExecute      = 0x17
	ComStmtSendLongData = 0x18
	ComStmtClose        = 0x19
	ComStmtReset        = 0x1a
	ComSetOption        = 0x1b
	ComStmtFetch        = 0x1c
	ComDaemon           = 0x1d
	ComBinlogDumpGtid   = 0x1e
	ComResetConnection  = 0x1f
)

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

type Packet struct {
	Length uint32
	SeqId  uint8
	Body   []byte
}

func NewPacket(seqId uint8, body []byte) *Packet {
	p := new(Packet)
	p.Length = uint32(len(body))
	p.SeqId = seqId
	p.Body = body
	return p
}

func (p Packet) ToBytes() []byte {
	var buf []byte
	buf = utils.AppendBuf(buf, utils.Uint24ToBytes(p.Length))
	buf = utils.AppendBuf(buf, utils.Uint8ToBytes(p.SeqId))
	buf = utils.AppendBuf(buf, p.Body)
	return buf
}

func HandShake() []byte {
	var buf []byte
	sequenceId := uint8(0)

	buf = append(buf)

	//protocal version
	buf = append(buf, 0x0a)

	// server-version
	buf = append(buf, []byte("muker-mysql-proxy-1.0")...)
	buf = append(buf, 0x00)

	// ConnectionID
	buf = append(buf, Uint32ToBytes(10001)...)

	// auth-plugin-data-part-1
	buf = append(buf, []byte("12345678")...)

	// filter
	buf = append(buf, 0x00)

	// capability flags
	buf = append(buf, 0xf7)
	buf = append(buf, 0xff)

	// character set
	buf = append(buf, 0x08)

	// status flags
	buf = append(buf, 0x00)
	buf = append(buf, 0x02)

	// Unused bytes 13bytes
	unused := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	buf = append(buf, unused[:]...)

	// auth-plugin-data-part-2
	for i := 0; i < 13; i++ {
		buf = append(buf, 0x00)
	}
	return NewPacket(sequenceId, buf).ToBytes()

}

func OK(sequenceId uint8) []byte {
	var buf []byte

	//header
	buf = append(buf, 0x00)

	// affected rows
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)

	//last_insert_id
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)

	//server status
	buf = append(buf, 0x00)
	buf = append(buf, 0x02)

	//warning
	buf = append(buf, 0x00)
	buf = append(buf, 0x00)

	return NewPacket(sequenceId, buf).ToBytes()
}

func ERR() []byte {
	var buf []byte
	sequenceId := uint8(0)

	//header
	buf = append(buf, 0xff)

	//error code
	buf = append(buf, 034)
	buf = append(buf, 012)

	//
	buf = append(buf, []byte("Internal error found")...)

	return NewPacket(sequenceId, buf).ToBytes()
}

func Quit() []byte {
	buf := []byte{1}
	return NewPacket(0, buf).ToBytes()
}
