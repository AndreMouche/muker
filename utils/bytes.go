package utils

func Uint24ToBytes(v uint32) []byte {
	v = v & 0x00ffffff
	b := make([]byte, 3)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	return b
}

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

func Uint8ToBytes(v uint8) []byte {
	b := make([]byte, 1)
	b[0] = v
	return b
}

func AppendBuf(buf []byte, toAppendBuf []byte) []byte {
	return append(buf, toAppendBuf[:]...)
}

func BytesToUint24(buf []byte) int {
	if len(buf) < 3 {
		return 0
	}
	v := int(buf[0])
	a := buf[1] << 8
	b := buf[2] << 16
	v += int(a)
	v += int(b)
	return v
}

func BytesToUint8(buf []byte) uint8 {
	return buf[0]
}
