package network

import (
	"sync/atomic"
)

type SafeIncrementer32 struct {
	val int32
}

func NewSafeIncrementer() *SafeIncrementer32 {
	return &SafeIncrementer32{
		val: 0,
	}
}

func (inc *SafeIncrementer32) Increment() int32 {
	for {
		old := inc.Value()
		new := old + 1
		if atomic.CompareAndSwapInt32(&inc.val, old, new) {
			return new
		}
	}
}

func (inc *SafeIncrementer32) Value() int32 {
	return atomic.LoadInt32(&inc.val)
}
