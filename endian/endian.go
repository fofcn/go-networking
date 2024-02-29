package main

import (
	"fmt"
)

func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
		fmt.Printf("%08b\n", buf)
	}
	buf[i] = byte(x)

	fmt.Printf("%08b\n", buf)
	return i + 1
}
