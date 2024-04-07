package once_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnce_ShouldExecuteOnce_WhenExecuteOnlyOnce(t *testing.T) {
	var i int
	once := sync.Once{}
	once.Do(func() {
		i = 1000
	})

	assert.Equal(t, 1000, i)
}

func TestOnce_ShouldExecuteOnce_WhenExecutedOnceAndExecuteAgain(t *testing.T) {
	var i int
	once := sync.Once{}

	once.Do(func() {
		i = 1000
	})

	once.Do(func() {
		i = 2000
	})

	assert.Equal(t, 1000, i)
}

func TestOnce_ShouldNotExecute_WhenFunctionPanicked(t *testing.T) {
	var num int
	once := sync.Once{}

	// 第一次调用应该 panic
	assert.Panics(t, func() {
		once.Do(func() {
			panic("Error")
		})
	})

	// 因为第一次panic了，此时再调用 Do 方法，函数 f 不应该被执行
	once.Do(func() {
		num = 1000
	})
	// 因为 num 的值没有被改变，所以应该还是 0
	assert.Equal(t, num, 0)
}

func TestOnce_ShouldExecuteAgain_WhenPreviousExecutionPaniced(t *testing.T) {
	var num int = 0
	once := sync.Once{}

	// 第一次调用应该 panic
	assert.Panics(t, func() {
		once.Do(func() {
			panic("Error")
		})
	})

	// 第一次 panic 后，下一次调用应该正确执行
	once.Do(func() {
		num = 1000
	})

	assert.Equal(t, 0, num)
}

func TestOnce_ShouldExecuteOnce_WhenCalledInMultipleGoroutines(t *testing.T) {
	var num atomic.Int32
	once := sync.Once{}
	wg := sync.WaitGroup{}
	wg.Add(100)

	// 启动 100 个 goroutine, 都尝试执行once.Do
	for i := 0; i < 100; i++ {
		go func() {
			once.Do(func() {
				num.Add(1)
			})
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(1), num.Load())
}
