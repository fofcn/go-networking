package main

import "fmt"

type EagerInstance struct {
}

var eagerInstance = &EagerInstance{}

func GetEagerInstance() *EagerInstance {
	return eagerInstance
}

func main() {

	s3 := GetEagerInstance()
	s4 := GetEagerInstance()
	if s3 == s4 {
		fmt.Println("s3 == s4")
	} else {
		fmt.Println("s3 != s4")
	}
}
