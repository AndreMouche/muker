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

const (
	PACKET_HEAD_LEN = 4
)

var ComSupported = make(map[byte]bool)

func init() {
	// from ComRefresh to ComResetConnection
	for i := 0x00; i < 0x1f; i++ {
		ComSupported[byte(i)] = false
	}
	ComSupported[ComQuit] = true
	ComSupported[ComInitDB] = true
	ComSupported[ComQuery] = true
	ComSupported[ComFieldList] = true
	ComSupported[ComCreateDB] = true
	ComSupported[ComDropDB] = true
}

type Packet struct {
	PktLen     uint32
	SequenceId uint8
	Buf        []byte
}

func NewPacket(pktLen uint32, sequenceId uint8, buf []byte) *Packet {
	return &Packet{
		PktLen:     pktLen,
		SequenceId: sequenceId,
		Buf:        buf,
	}
}

func WritePacket(buf []byte, sequenceId uint8) ([]byte, error) {
	pktLen := len(buf)
	ret := make([]byte, 4+pktLen)
	pktLenBuf := utils.Uint24ToBytes(uint32(pktLen))
	ret[0] = pktLenBuf[0]
	ret[1] = pktLenBuf[1]
	ret[2] = pktLenBuf[2]
	ret[3] = byte(sequenceId)
	copy(ret[4:], buf)
	return ret, nil
}

type HandShakePacket struct {
	protoVersion    byte // 1 byte
	serverVersion   string
	connId          uint32 // 4 byte
	authPluginData0 []byte // 8 byte
	filter          byte   // 1 byte
	capabilityFlags []byte // 2 byte
	characterSet    byte   // 1 byte
	statusFlags     []byte // 2 byte
	unusedByte      []byte // 13 byte
	authPluginData1 []byte // 13 byte
}

func DefaultHandShakePacket(connId uint32) *HandShakePacket {
	return &HandShakePacket{
		protoVersion:    0x0a,
		serverVersion:   "Muker-MySQL-Proxy-1.0",
		connId:          connId,
		authPluginData0: utils.RandBytes(8),
		filter:          byte(0x00),
		capabilityFlags: []byte{0xf7, 0xff},
		characterSet:    byte(0x08),
		statusFlags:     []byte{0x00, 0x02},
		unusedByte:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		authPluginData1: append(utils.RandBytes(12), byte(0x00)),
	}
}

// TODO optimizition.
func (p *HandShakePacket) Write(sequenceId uint8) ([]byte, error) {
	idx := 0
	buf := make([]byte, p.Len())

	buf[0] = p.protoVersion
	idx += 1

	copy(buf[idx:], []byte(p.serverVersion))
	idx += len(p.serverVersion)

	buf[idx] = byte(0x00)
	idx += 1

	copy(buf[idx:], utils.Uint32ToBytes(p.connId))
	idx += 4

	copy(buf[idx:], p.authPluginData0)
	idx += len(p.authPluginData0)

	buf[idx] = p.filter
	idx += 1

	copy(buf[idx:], p.capabilityFlags)
	idx += len(p.capabilityFlags)

	buf[idx] = p.characterSet
	idx += 1

	copy(buf[idx:], p.statusFlags)
	idx += len(p.statusFlags)

	copy(buf[idx:], p.unusedByte)
	idx += len(p.unusedByte)

	copy(buf[idx:], p.authPluginData1)
	idx += len(p.authPluginData1)

	return WritePacket(buf, sequenceId)
}

func (p *HandShakePacket) Read(buf []byte) error {
	return nil
}

func (p *HandShakePacket) Len() int {
	return 1 + len(p.serverVersion) + 1 + 4 + 8 + 1 + 2 + 1 + 2 + 13 + 13
}

type OkPacket struct {
	header        byte     // 1 byte
	affetctedRows uint64   // lenenc
	lastInsertId  uint64   // lenenc
	statusFlags   []byte   // 2 byte
	warnings      []byte   // 2 byte
	warnMsgs      []string // <lenenc>string-bytes
}

func DefaultOkPacket() *OkPacket {
	return &OkPacket{
		header:        0x00,
		affetctedRows: 0,
		lastInsertId:  0,
		statusFlags:   []byte{0x00, 0x02},
		warnings:      []byte{0x00, 0x00},
	}
}

func (p *OkPacket) Write(sequenceId uint8) ([]byte, error) {
	idx := 0
	buf := make([]byte, p.Len())

	buf[0] = p.header
	idx += 1

	affectedRowBuf, _ := utils.IntToLenEncode(p.affetctedRows)
	copy(buf[idx:], affectedRowBuf)
	idx += len(affectedRowBuf)

	lastInsertIdBuf, _ := utils.IntToLenEncode(p.lastInsertId)
	copy(buf[idx:], lastInsertIdBuf)
	idx += len(lastInsertIdBuf)

	copy(buf[idx:], p.statusFlags)
	idx += 2

	if len(p.warnings) > 0 {
		for _, warnMsg := range p.warnMsgs {
			warnLenBuf, _ := utils.IntToLenEncode(uint64(len(warnMsg)))
			copy(buf[idx:], warnLenBuf)
			idx += len(warnLenBuf)

			warnBuf := []byte(warnMsg)
			copy(buf[idx:], warnBuf)
			idx += len(warnBuf)
		}
	}
	return WritePacket(buf, sequenceId)
}

func (p *OkPacket) Read(buf []byte) error {
	return nil
}

func (p *OkPacket) Len() int {
	ret := 0
	ret += 1

	affectedRowsBuf, _ := utils.IntToLenEncode(p.affetctedRows)
	ret += len(affectedRowsBuf)

	lastInsertIdBuf, _ := utils.IntToLenEncode(p.lastInsertId)
	ret += len(lastInsertIdBuf)

	// statusFlags
	ret += 2

	// warnings
	ret += 2

	for _, s := range p.warnMsgs {
		size := len(s)
		strBuf, _ := utils.IntToLenEncode(uint64(size))
		ret += len(strBuf)
		ret += size
	}
	return ret
}

type ErrorPacket struct {
	header         byte   // 1 byte
	errorCode      uint16 // 2 byte
	sqlStateMarker byte   // 1 byte
	sqlState       []byte // 5 byte
	errMsg         string // string<EOF>
}

func DefaultErrorPacket(errStr string) *ErrorPacket {
	return &ErrorPacket{
		header:         0xff,
		errorCode:      1024,
		sqlStateMarker: byte(0x23),
		sqlState:       []byte{0x48, 0x59, 0x30, 0x30, 0x30},
		errMsg:         errStr,
	}
}

func (p *ErrorPacket) Write(sequenceId uint8) ([]byte, error) {

	idx := 0
	buf := make([]byte, p.Len())

	buf[idx] = p.header
	idx += 1

	errorCodeBuf := utils.Uint16ToBytes(p.errorCode)
	copy(buf[idx:], errorCodeBuf)
	idx += 2

	buf[idx] = p.sqlStateMarker
	idx += 1

	copy(buf[idx:], p.sqlState)
	idx += 5

	if len(p.errMsg) > 0 {
		errMsgBuf := []byte(p.errMsg)
		copy(buf[idx:], errMsgBuf)
		idx += len(errMsgBuf)
	}

	return WritePacket(buf, sequenceId)
}

func (p *ErrorPacket) Read(buf []byte) error {
	return nil
}

func (p *ErrorPacket) Len() int {
	return 1 + 2 + 1 + 5 + len(p.errMsg)
}

type CommandPacket struct {
	comType byte   // 1 byte
	query   string // string<EOF>
}

func DefaultCommandPacket() *CommandPacket {
	return &CommandPacket{
		comType: ComQuery,
		query:   "select @@version",
	}
}

func (p *CommandPacket) Write(sequenceId uint8) ([]byte, error) {

	idx := 0
	buf := make([]byte, p.Len())

	buf[0] = p.comType
	idx += 1

	if len(p.query) > 0 {
		queryBuf := []byte(p.query)
		copy(buf[idx:], queryBuf)
		idx += len(queryBuf)
	}

	return WritePacket(buf, sequenceId)
}

func (p *CommandPacket) Read(buf []byte) error {
	return nil
}

func (p *CommandPacket) Len() int {
	return 1 + len(p.query)
}

type ColumnDefPacket struct {
}

func (p *ColumnDefPacket) Write(sequenceId uint8) ([]byte, error) {
}

func (p *ColumnDefPacket) Read(buf []byte) error {
}

func (p *ColumnDefPacket) Len() int {
}

type ResultSetRowPacket struct {
	fieldCount int
	fields     []string
}

func DefaultResultSetRowPacket(fieldCount int) *ResultSetRowPacket {
	return &ResultSetRowPacket{
		fieldCount: fieldCount,
		fields:     make([]string, fieldCount),
	}
}

func (p *ResultSetRowPacket) Write(sequenceId uint8) ([]byte, error) {
}

func (p *ResultSetRowPacket) Read(buf []byte) error {
	// NULL is return
	if buf[0] == 0xfb {
		p.fieldCount = 0
		p.fields = nil
		return nil
	}
	fieldIdx := 0
	idx := 0
	bufLen := len(buf)
	for {
		if iLen, err := utils.LenEncodeToInt(buf[idx:]); err != nil {
			return err
		}
		if calcBytes, err2 := uilts.CalcLenForLenEncode(buf[idx:]); err != nil {
			return err
		}
		idx += calcBytes
		if idx+iLen > bufLen {
			return fmt.Errorf("ResultSetRowPacket read overflow: bufLen: %d, readIndex: %d", bufLen, idx+iLen)
		}
		fields[fieldIdx] = string(buf[idx : idx+iLen])
		idx += iLen
		fieldIdx++
		if fields == fieldIdx {
			if idx != bufLen {
				return fmt.Errorf("ResultSetRowPacket unfinished bytes.")
			}
			break
		}
	}
	return nil
}

func (p *ResultSetRowPacket) Len() int {
	if p.fieldCount == 0 {
		return 1
	}
	ret := 0
	for i := 0; i < p.fieldCount; i++ {
		if strLen, err := utils.IntToLenEncode(len(p.fields[i])); err != nil {
			return 0
		}
		ret += strLen
	}
	return ret
}
