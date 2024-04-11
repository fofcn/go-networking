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

type EagerInstance struct {
}

var eagerInstance = &EagerInstance{}

func GetEagerInstance() *EagerInstance {
	return eagerInstance
}

func main() {
	s1 := GetLazyInstance()
	s2 := GetLazyInstance()
	if s1 == s2 {
		fmt.Println("s1 == s2")
	} else {
		fmt.Println("s1 != s2")
	}

	s3 := GetEagerInstance()
	s4 := GetEagerInstance()
	if s3 == s4 {
		fmt.Println("s3 == s4")
	} else {
		fmt.Println("s3 != s4")
	}
}
