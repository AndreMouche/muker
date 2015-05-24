package utils

import (
	"math/rand"
)

func RandBytes(size int) []byte {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(rand.Int31n(128))
	}
	return buf
}
