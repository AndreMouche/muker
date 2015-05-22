package proto

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

func HandShake() []byte {
	var buf []byte

	//protocal version
	buf = append(buf, 0x0a)

	// server-version
	buf = append(buf, []byte("Muker -- MySQL Proxy")...)
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

	//

	return buf
}
