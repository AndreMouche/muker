package utils

import (
	"fmt"
)

func Uint24ToBytes(v uint32) []byte {
	v = v & 0x00ffffff
	b := make([]byte, 3)
	b[0] = byte(v)
	b[1] = byte((v >> 8) & 0xff)
	b[2] = byte((v >> 16) & 0xff)
	return b
}

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte((v >> 8) & 0xff)
	b[2] = byte((v >> 16) & 0xff)
	b[3] = byte((v >> 24) & 0xff)
	return b
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(v)
	b[1] = byte((v >> 8) & 0xff)
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

func IntToLenEncode(v uint64) ([]byte, error) {
	var buf []byte
	if v < 251 {
		buf = append(buf, byte(v&0xff))
	} else if v < (1 << 16) {
		buf = append(buf, byte(0xfc))
		buf = append(buf, byte(v&0xff))
		buf = append(buf, byte((v>>8)&0xff))
	} else if v < (1 << 24) {
		buf = append(buf, byte(0xfd))
		buf = append(buf, byte(v&0xff))
		buf = append(buf, byte((v>>8)&0xff))
		buf = append(buf, byte((v>>16)&0xff))
	} else {
		buf = append(buf, byte(0xfe))
		buf = append(buf, byte(v&0xff))
		buf = append(buf, byte((v>>8)&0xff))
		buf = append(buf, byte((v>>16)&0xff))
		buf = append(buf, byte((v>>24)&0xff))
		buf = append(buf, byte((v>>32)&0xff))
		buf = append(buf, byte((v>>40)&0xff))
		buf = append(buf, byte((v>>48)&0xff))
		buf = append(buf, byte((v>>56)&0xff))
	}
	return buf, nil
}

func LenEncodeToInt(buf []byte) (uint64, error) {
	size := len(buf)
	if size < 1 {
		return 0, fmt.Errorf("LenEncodeToInt size < 1")
	}
	if buf[0] < 0xfb {
		return uint64(buf[0]), nil
	}
	v := uint64(0)
	if buf[0] == 0xfc {
		if size < 3 {
			return 0, fmt.Errorf("LenEncodeToInt size < 3")
		}
		v |= uint64(buf[1])
		v |= uint64(buf[2]) << 8
	} else if buf[0] == 0xfd {
		if size < 4 {
			return 0, fmt.Errorf("LenEncodeToInt size < 4")
		}
		v |= uint64(buf[1])
		v |= uint64(buf[2]) << 8
		v |= uint64(buf[3]) << 16
	} else if buf[0] == 0xfe {
		if size < 9 {
			return 0, fmt.Errorf("LenEncodeToInt size < 9")
		}
		v |= uint64(buf[1])
		v |= uint64(buf[2]) << 8
		v |= uint64(buf[3]) << 16
		v |= uint64(buf[4]) << 24
		v |= uint64(buf[5]) << 32
		v |= uint64(buf[6]) << 40
		v |= uint64(buf[7]) << 48
		v |= uint64(buf[8]) << 56
	}
	return v, nil
}

func CalcLenForLenEncode(buf []byte) (uint64, error) {
	if buf[0] < 0xfb {
		return 1, nil
	} else if buf[0] == 0xfc {
		return 3, nil
	} else if buf[0] == 0xfd {
		return 4, nil
	} else if buf[0] == 0xfe {
		return 9, nil
	} else {
		return 0, fmt.Errorf("CalcLenForLenEncode failed. %x", buf)
	}
}
