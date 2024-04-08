package mutex_test

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMutex_ShouldLock_WhenNotLockedBefore 测试互斥锁在未锁定状态时能否成功锁定。
// 参数 t 用于测试上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldLock_WhenNotLockedBefore(t *testing.T) {
	// 创建一个互斥锁
	mux := sync.Mutex{}
	// 尝试锁定互斥锁
	mux.Lock()
	// 标记锁定是否成功
	isSuccess := true
	// 确保在函数退出时解锁，避免死锁
	defer mux.Unlock()

	// 验证锁定是否成功
	assert.True(t, isSuccess)
}

// TestMutex_ShouldBlock_WhenUsingLockAndOneRoutineHasLocked 测试当一个协程已经锁定互斥锁时，
// 其他协程尝试锁定是否会被阻塞。
// 参数 t 用于测试上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldBlock_WhenUsingLockAndOneRoutineHasLocked(t *testing.T) {
	// 创建一个互斥锁
	mux := sync.Mutex{}
	// 创建一个等待组，用于同步协程
	wg := sync.WaitGroup{}
	// 添加一个计数，表示要等待的一个协程
	wg.Add(1)

	// 启动一个协程去锁定互斥锁
	go func() {
		// 表示协程启动完成
		wg.Done()
		// 尝试锁定互斥锁
		mux.Lock()
		// 确保在函数退出时解锁
		defer mux.Unlock()
		// 持锁一段时间，模拟锁定状态
		time.Sleep(5 * time.Second)
	}()

	// 等待协程启动完成
	wg.Wait()
	println("go routine started")
	// 主协程尝试锁定互斥锁，应被阻塞
	mux.Lock()
	// 标记锁定是否成功
	isSuccess := true
	// 确保解锁
	defer mux.Unlock()

	// 验证主协程是否成功锁定
	assert.True(t, isSuccess)
}

// TestMutex_ShouldLocked_WhenTryLockAndNoOtherLockers 测试在没有其他锁定的情况下，
// 尝试使用 TryLock 方法是否能成功锁定互斥锁。
// 参数 t 用于测试上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldLocked_WhenTryLockAndNoOtherLockers(t *testing.T) {
	// 创建一个互斥锁
	mux := sync.Mutex{}
	// 尝试立即锁定互斥锁
	isSuccess := mux.TryLock()
	// 确保解锁
	defer mux.Unlock()

	// 验证是否成功锁定
	assert.True(t, isSuccess)
}

// TestMutex_ShouldNotLocked_WhenTryLockAndOtherLockers 测试当已有协程锁定互斥锁时，
// 尝试使用 TryLock 方法是否不能成功锁定互斥锁。
// 参数 t 用于测试上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldNotLocked_WhenTryLockAndOtherLockers(t *testing.T) {
	// 创建一个互斥锁
	mux := sync.Mutex{}
	// 创建一个等待组，用于同步协程
	wg := sync.WaitGroup{}
	// 添加一个计数，表示要等待的一个协程
	wg.Add(1)
	// 启动一个协程去锁定互斥锁
	go func() {
		// 表示协程启动完成
		//这里先调用Done是为了不继续阻塞wg.Wait()让wg.Wait()之后的mutex.Lock可以继续执行
		wg.Done()
		// 尝试锁定互斥锁
		mux.Lock()
		// 确保在函数退出时解锁
		defer mux.Unlock()
		// 持锁一段时间，模拟锁定状态
		time.Sleep(5 * time.Second)
	}()

	// 等待协程启动完成
	wg.Wait()
	// 尝试立即锁定互斥锁
	isSuccess := mux.TryLock()
	// 确保解锁
	defer mux.Unlock()

	// 验证是否未能成功锁定
	assert.False(t, isSuccess)
}

// TestMutex_ShouldIncrementCounterSuccess_WhenUseMultipleGoroutineAndAddCounterInLock 测试在多个协程并发情况下，使用Mutex锁定来增加计数器是否成功
// 参数:
// - t *testing.T: 测试环境的句柄，用于报告测试失败和日志记录
// 返回值: 无
func TestMutex_ShouldIncrementCounterSuccess_WhenUseMultipleGoroutineAndAddCounterInLock(t *testing.T) {
	// 初始化互斥锁、等待组和计数器
	mux := sync.Mutex{}
	wg := sync.WaitGroup{}
	counter := 0
	// 设置运行时的最大协程数为10，以模拟并发环境
	runtime.GOMAXPROCS(10)

	// 添加100个协程来并发执行增加计数器的操作
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			// 在协程结束时释放等待组
			defer wg.Done()
			// 加锁以确保对计数器的操作是互斥的
			mux.Lock()
			defer mux.Unlock()
			// 增加计数器
			counter++
		}()
	}
	// 等待所有协程完成
	wg.Wait()

	// 验证计数器是否被正确地增加了100次
	assert.Equal(t, 100, counter)
}

// TestMutex_ShouldPanic_WhenUnlockWithoutLock 测试当互斥锁没有被锁定时，尝试解锁是否会引发Panic。
// 参数 t *testing.T 用于测试的上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldPanic_WhenUnlockWithoutLock(t *testing.T) {
	// 创建一个 sync.Mutex 实例。
	mux := sync.Mutex{}
	// 使用 assert 包的 Panics 函数来断言是否会引发Panic。
	assert.Panics(t, func() {
		// 尝试解锁一个没有被锁定的互斥锁。
		mux.Unlock()
	})
}

// TestMutex_ShouldStarvation_WhenOneRoutineHoldLockedForLongTime 测试当一个协程长时间持有互斥锁时，
// 其他协程是否会出现饥饿现象。参数 t *testing.T 用于测试的上下文，提供测试控制和日志记录功能。
func TestMutex_ShouldStarvation_WhenOneRoutineHoldLockedForLongTime(t *testing.T) {
	// 定义一个 sync.Mutex 实例用于测试互斥锁饥饿问题。
	var mu sync.Mutex
	// 用于同步协程开始的 waitgroup。
	var start sync.WaitGroup
	// 用于同步协程结束的 waitgroup。
	var done sync.WaitGroup

	start.Add(1) // 准备启动一个协程。
	done.Add(1)  // 等待一个协程完成。
	go func() {  // 启动一个协程长时间持有锁。
		start.Done() // 表示协程已启动并准备好。
		mu.Lock()    // 获取锁并长时间持有。
		time.Sleep(1000 * time.Second)
		mu.Unlock() // 最终释放锁。
		done.Done() // 表示协程已完成。
	}()

	start.Wait()                 // 等待协程开始并持有锁。
	time.Sleep(time.Millisecond) // 稍微延时以确保锁被持有。

	start.Add(1) // 准备启动另一个协程。
	done.Add(1)  // 等待另一个协程完成。
	go func() {  // 启动另一个协程尝试获取锁，以测试是否出现饥饿。
		start.Done()                             // 表示协程已启动并准备好。
		mu.Lock()                                // 尝试获取锁。
		t.Log("Starving goroutine got the lock") // 如果获取到锁，则记录日志。
		mu.Unlock()                              // 最终释放锁。
		done.Done()                              // 表示协程已完成。
	}()

	start.Wait()                 // 等待第二个协程开始。
	time.Sleep(time.Millisecond) // 稍微延时以确保尝试获取锁的协程已运行。

	mu.Lock()                            // 主协程尝试获取锁，以进一步测试饥饿情况。
	t.Log("Main goroutine got the lock") // 如果获取到锁，则记录日志。
	mu.Unlock()                          // 释放主协程持有的锁。

	done.Wait() // 等待所有协程完成，确保测试完整执行。
}
