package main

import (
	"fmt"
	"unsafe"
)

func createSlice(cap int) []int {
	return make([]int, cap)
}

func main() {

	// s := make([]int, 20)
	// s = append(s, 1)
	// println(len(s))

	s := make([]int, 5)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	s[3] = 4
	s[4] = 5
	s = append(s, 1, 2, 3, 4, 5)
	fmt.Printf("s slice addr: %p\n", &s)
	println("s addr: ", unsafe.Pointer(&s[0]), unsafe.Pointer(&s[1]), unsafe.Pointer(&s[2]), unsafe.Pointer(&s[3]), unsafe.Pointer(&s[4]))

	// s1 := s[:1]
	// fmt.Printf("s1 slice addr: %p\n", &s1)
	// println("s1 addr: ", unsafe.Pointer(&s1[0]))

	// s2 := append(s[0:2], s[3:]...)
	// fmt.Printf("s2 slice addr: %p\n", &s2)
	// println("s2 addr: ", unsafe.Pointer(&s2[0]))

	scopy := make([]int, 10)
	copy(s, scopy)
	println("scopy 1st element addr: ", unsafe.Pointer(&scopy[0]))

	// int64slice := make([]int, 1<<32)
	// copy(s, int64slice)

}
