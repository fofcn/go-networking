package main

import "fmt"

var lazyinstance *LazySingleton

type LazySingleton struct {
}

func GetLazyInstance() *LazySingleton {
	if lazyinstance == nil {
		lazyinstance = &LazySingleton{}
	}
	return lazyinstance
}

func main() {
	s1 := GetLazyInstance()
	s2 := GetLazyInstance()
	if s1 == s2 {
		fmt.Println("s1 == s2")
	} else {
		fmt.Println("s1 != s2")
	}
}
