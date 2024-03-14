package network_test

import (
	"go-networking/network"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementWhenCalledConcurrentlyShouldReturnCorrectValue(t *testing.T) {
	incrementer := network.NewSafeIncrementer()
	n := 1000
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			incrementer.Increment()
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(n), incrementer.Value())
}

func TestValueWhenCalledConcurrentlyShouldReturnSameValue(t *testing.T) {
	incrementer := network.NewSafeIncrementer()
	var wg sync.WaitGroup
	n := 1000

	// 先做 n 次增量操作，让 incrementer 的值变为 n
	for i := 0; i < n; i++ {
		incrementer.Increment()
	}

	// 并发地获取 incrementer 的值，并检查是否和预期一致
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if incrementer.Value() != int32(n) {
				t.Errorf("Expected value %d, but got %d", n, incrementer.Value())
			}
		}()
	}
	wg.Wait()
}
