package atomic

import (
	"sync/atomic"
)

type Bool struct {
	val atomic.Int32
}

func NewBool(b bool) *Bool {
	intVal := int32(0)
	if b {
		intVal = 1
	}
	new := &Bool{}
	new.val.Store(intVal)
	return new
}

// Get atomically retrieves the current boolean value.
func (b *Bool) Get() bool {
	return b.val.Load() != 0
}

// Set atomically sets the boolean value.
func (b *Bool) Set(val bool) {
	intVal := int32(0)
	if val {
		intVal = 1
	}
	b.val.Store(intVal)
}

// CompareAndSet atomically sets the value to the given updated value
// if the current value == the expected value.
func (b *Bool) CompareAndSet(expected, newValue bool) bool {
	var oldInt int32 = 0
	if expected {
		oldInt = 1
	}
	var newInt int32 = 0
	if newValue {
		newInt = 1
	}
	return b.val.CompareAndSwap(oldInt, newInt)
}
