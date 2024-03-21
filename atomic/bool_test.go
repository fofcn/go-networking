package atomic_test

import (
	"go-networking/atomic"
	"testing"
)

func BenchmarkAtomicBoolGet(b *testing.B) {
	ab := atomic.NewBool(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ab.Get()
	}
}

func BenchmarkAtomicBoolSet(b *testing.B) {
	ab := atomic.NewBool(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ab.Set(true)
	}
}

func BenchmarkAtomicBoolCompareAndSet(b *testing.B) {
	ab := atomic.NewBool(true)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ab.CompareAndSet(true, false)
	}
}
