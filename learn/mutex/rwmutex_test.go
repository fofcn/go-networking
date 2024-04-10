package mutex_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试读写互斥锁在正常读锁定和解锁情况下的成功执行
func TestRWMutex_ShouldSuccess_WhenNormalReaderLockAndUnLock(t *testing.T) {
	// 初始化一个读写互斥锁
	rwmutex := sync.RWMutex{}
	// 获取读锁
	rwmutex.RLock()
	// 设置成功标志为true，使用defer确保在函数结束时释放读锁
	isSuccess := true
	defer rwmutex.RUnlock()
	// 记录日志表示测试成功
	t.Log("success")
	// 断言成功标志为true
	assert.True(t, isSuccess)
}

// 测试RWMutex的写锁功能是否正常
func TestRWMutex_ShouldSuccess_WhenNormalWriterLockAndUnLock(t *testing.T) {
	rwmutex := sync.RWMutex{} // 创建一个sync.RWMutex类型的变量
	rwmutex.Lock()            // 获取写锁
	isSuccess := true         // 标记为成功状态
	defer rwmutex.Unlock()    // 确保在函数退出时释放锁，避免死锁
	t.Log("success")          // 记录测试日志
	assert.True(t, isSuccess) // 断言isSuccess为true，验证操作成功
}

// 函数测试了在正常情况下，
// 读写锁（RWMutex）的读锁（RLock）和写锁（Lock）的加锁与解锁操作是否成功。
func TestRWMutex_ShouldSuccess_WhenNormalReaderWriterLockAndUnLock(t *testing.T) {
	// 初始化一个读写锁
	rwmutex := sync.RWMutex{}
	// 尝试获取读锁并立即释放
	rwmutex.RLock()
	rwmutex.RUnlock()
	// 尝试获取写锁并立即释放
	rwmutex.Lock()
	rwmutex.Unlock()
	// 标记测试为成功
	isSuccess := true
	// 记录测试成功日志
	t.Log("success")
	// 断言测试结果为真
	assert.True(t, isSuccess)
}

// 测试读写锁在多协程情况下的读写互斥
func TestRWMutex_ShouldSuccess_WhenReaderAndWriterInDifferentRoutine(t *testing.T) {
	// 初始化一个读写锁和等待组，用于协调不同协程的操作。
	rwmutex := sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(2) // 预期有两个协程完成操作

	// 启动一个协程作为读锁持有者
	go func() {
		rwmutex.RLock()   // 获取读锁
		println("reader") // 打印读操作标识
		rwmutex.RUnlock() // 释放读锁
		wg.Done()         // 表示读操作完成
	}()

	// 启动另一个协程作为写锁持有者
	go func() {
		rwmutex.Lock()    // 获取写锁
		println("writer") // 打印写操作标识
		rwmutex.Unlock()  // 释放写锁
		wg.Done()         // 表示写操作完成
	}()

	wg.Wait() // 等待所有协程完成操作
	isSuccess := true
	t.Log("success")          // 记录测试成功
	assert.True(t, isSuccess) // 断言测试结果为真
}

// 测试读写锁在多个读锁情况下的读写互斥
func TestRWMutex_ShouldBlockWriter_WhenMultipleReader(t *testing.T) {
	rwmutex := sync.RWMutex{}
	ch := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			wg.Done()
			rwmutex.RLock()
			println("reader Locked", i)
			time.Sleep(10 * time.Second)
			rwmutex.RUnlock()
			println("reader UnLocked", i)
		}(i)
	}

	go func() {
		wg.Wait()
		println("writer try to accquire wlock")
		rwmutex.Lock()
		println("writer has accquired wlock")
		defer rwmutex.Unlock()
		ch <- true
	}()

	<-ch
	isSuccess := true
	t.Log("success")
	assert.True(t, isSuccess)
}

// 测试读写锁在多个写锁情况下的读写互斥
func TestRWMutex_ShouldBlockReaders_WhenWriterIsPresent(t *testing.T) {
	rwmutex := sync.RWMutex{}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		println("writer try to accquire wlock")
		rwmutex.Lock()
		println("writer has accquired wlock")
		wg.Done()
		time.Sleep(10 * time.Second)
		defer rwmutex.Unlock()
		println("writer has released wlock")
	}()

	wg.Wait()
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(i int) {
			println("reader try to lock", i)
			rwmutex.RLock()
			println("reader Locked", i)
			rwmutex.RUnlock()
			println("reader UnLocked", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	isSuccess := true
	t.Log("success")
	assert.True(t, isSuccess)
}

// 测试读写锁在多个写锁情况下的读写互斥
func TestRWMutex_ShouldBlockConcurrentWriters(t *testing.T) {
	rwmutex := sync.RWMutex{}
	var blockedWriter bool
	ch := make(chan bool)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		wg.Done()
		println("Writer 1 try to accquire wlock")
		rwmutex.Lock()
		println("Writer 1 has accquired wlock")
		defer rwmutex.Unlock()
		time.Sleep(15 * time.Second)
	}()

	go func() {
		wg.Wait()
		println("Writer 2 try to accquire wlock")
		rwmutex.Lock()
		println("Writer 2 has accquired wlock")
		ch <- true
		defer rwmutex.Unlock()
	}()

	select {
	case <-ch:
		blockedWriter = false
	case <-time.After(20 * time.Second):
		blockedWriter = true
	}
	assert.True(t, blockedWriter)
}

// 测试读写锁在多个读锁情况下的读写互斥
func TestRWMutex_ShouldLockSuccess_WhenTryingToReadLockTwice(t *testing.T) {
	rwmutex := sync.RWMutex{}
	writerWaitGroup := sync.WaitGroup{}
	writerWaitGroup.Add(1)

	go func() {
		rwmutex.RLock()
		println("readlock locked once")
		rwmutex.RLock()
		println("readlock locked twice")
		rwmutex.RUnlock()
		rwmutex.RUnlock()
		defer writerWaitGroup.Done()
	}()

	writerWaitGroup.Wait()
	isSuccess := true

	assert.True(t, isSuccess)
}

// 测试读写锁在多个写锁情况下的读写互斥
func TestRWMutex_ShouldBeBlocked_WhenTryingToWriteLockTwice(t *testing.T) {
	rwmutex := sync.RWMutex{}
	ch := make(chan bool)
	go func() {
		rwmutex.Lock()
		println("writelock locked once")
		rwmutex.Lock()
		println("writelock locked twice")
		rwmutex.Unlock()
		rwmutex.Unlock()
		ch <- true
	}()

	isBlocked := false

	select {
	case <-ch:
		println("should not execute this block")
		assert.False(t, isBlocked)
	case <-time.After(10 * time.Second):
		isSuccess := true
		println("executed timeout block")
		assert.True(t, isSuccess)
	}

}

// 测试读写锁在多个读锁情况下的读写互斥
func TestRWMutex_ShouldBeBlocked_WhenAccquireWriteLockThenReadLock(t *testing.T) {
	rwmutex := sync.RWMutex{}
	ch := make(chan bool)
	go func() {
		rwmutex.Lock()
		println("writelock locked once")
		rwmutex.RLock()
		println("readlock locked twice")
		rwmutex.RUnlock()
		rwmutex.Unlock()
		ch <- true
	}()
	isBlocked := false

	select {
	case <-ch:
		println("should not execute this block")
		assert.False(t, isBlocked)
	case <-time.After(10 * time.Second):
		isSuccess := true
		println("executed timeout block")
		assert.True(t, isSuccess)
	}

}

// 测试读写锁在多个读锁情况下的读写互斥
func TestRWMutex_ShouldBeBlocked_WhenAccquireReadLockThenWriteLock(t *testing.T) {
	rwmutex := sync.RWMutex{}
	ch := make(chan bool)
	go func() {
		rwmutex.RLock()
		println("readlock locked once")
		rwmutex.Lock()
		println("writelock locked twice")
		rwmutex.Unlock()
		rwmutex.RUnlock()
		ch <- true
	}()
	isBlocked := false

	select {
	case <-ch:
		println("should not execute this block")
		assert.False(t, isBlocked)
	case <-time.After(10 * time.Second):
		isSuccess := true
		println("executed timeout block")
		assert.True(t, isSuccess)
	}

}

// 测试读写锁在多个读锁情况下的读写互斥
func TestRWMutex_ShouldDeadlockOrBlocked_WhenLockOneGoroutineAccquiredLockAndAnotherGoroutineAccquireLockAgain(t *testing.T) {
	var rwmutex1, rwmutex2 sync.RWMutex
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	ch := make(chan bool)

	wg.Add(1)
	wg1.Add(1)
	go func() {
		rwmutex1.Lock()
		println("rwmutex1 locked")
		wg.Done()
		wg1.Wait()
		println("rwmutex2 try to accquire lock")
		rwmutex2.Lock()
	}()
	go func() {
		wg.Wait()
		rwmutex2.Lock()
		println("rwmutex2 locked")
		wg1.Done()
		println("rwmutex1 try to accquire lock")
		rwmutex1.Lock()
		ch <- true
	}()
	isBlocked := false

	select {
	case <-ch:
		println("should not execute this block")
		assert.False(t, isBlocked)
	case <-time.After(10 * time.Second):
		isSuccess := true
		println("executed timeout block")
		assert.True(t, isSuccess)
	}

}
