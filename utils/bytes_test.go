package utils

import (
	"testing"
)

func TestIntToLenEncode(t *testing.T) {
	// 9223372036854775807 = 1^63 - 1
	arr := []uint64{0, 250, 2333, 100000, 2342342342342, 234213345, 9223372036854775807}
	for _, a := range arr {
		buf, err := IntToLenEncode(a)
		if err != nil || len(buf) == 0 {
			t.Errorf("IntToLenEncode %d failed.", a)
		}
		v, err2 := LenEncodeToInt(buf)
		if err2 != nil || v != a {
			t.Errorf("LenEncodeToInt  %x to %d failed. origin: %d", buf, v, a)
		}
	}
}
