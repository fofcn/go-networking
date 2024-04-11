package main

import "sync"

type LockInstance struct {
	s string
}

var mux sync.Mutex
var lockinstance *LockInstance

func GetLockInstance() *LockInstance {
	if lockinstance == nil {
		mux.Lock()
		defer mux.Unlock()
		println("create singleton instance")
		if lockinstance == nil {
			lockinstance = &LockInstance{
				s: "lock",
			}
		}
	}

	return lockinstance
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance := GetLockInstance()
			println(instance.s)
		}()
	}
	wg.Wait()
}
