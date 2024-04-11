package main

type Singleton struct {
	s string
}

var instance *Singleton

func init() {
	instance = &Singleton{s: "hello"}
}

func GetInstance() *Singleton {
	return instance
}

func main() {
	println(GetInstance().s)
}
