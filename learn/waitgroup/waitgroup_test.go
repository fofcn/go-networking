package waitgroup_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 这是对Go语标准库中sync包下的WaitGroup的描述。

// WaitGroup用于等待一组并发的goroutine结结束。
// 主goroutine会调用Add方法来设置需要等待的goroutine的数量。
// 然后，每个goroutine在结束之后需要调用Done方法。
// 同时，可以使用Wait方法阻塞，直到所有的goroutine都结束。

// 在使用WaitGroup的过程中，一旦开始使用就不能再进行复制。

// 在 Go 内存模型的术语中，对 Done 的调用 “在” 它解除的任何 Wait 调用返回前“同步”。
// 这是一种保证内存可见性的机制，即一旦Done被调用，Wait就能获知并返回。
// 在并发编程中，这种机制能确保相互依赖的操作的正确顺序，避免了潜在的竞态条件。

// [很重要]如果你想复用一个WaitGroup，等待几组不相关的事件，你必须确保在调用新一轮的Add之前，所有先前的Wait调用都已经返回完毕。

// TestWaitGroup_ShouldComplete_WhenAddCountEqualsDoneCount 测试了一个正常的
// WaitGroup 使用场景。我们添加了一个需要等待的 goroutine，当这个 goroutine 完成时，
// WaitGroup 的 Wait 方法应该返回。这是最基础的使用 WaitGroup 的场景。
func TestWaitGroup_ShouldComplete_WhenAddCountEqualsDoneCount(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		println("Hello WaitGroup")
		wg.Done()
	}()

	wg.Wait()

	completed := 1
	assert.Equal(t, 1, completed)
}

// TestWaitGroup_ShouldCompleted_WhenMultipleGoroutineDoNormalAddAndDone 用于测试多个
// 并发的 Goroutine 是否能够正确地被 WaitGroup 跟踪以及同步。
// 我们创建了 loopCount 个并发的 Goroutine，并且每个 Goroutine 会在完成任务后调用 Done 方法标记。
// 我们在主 Goroutine 中 阻塞等待所有的 Goroutine 完成。
// 如果所有的 goroutine 都能顺利地完成数据处理且调用了Done 方法，WaitGroup 的 Wait 方法将会顺利返回，解除阻塞并且继续执行。
func TestWaitGroup_ShouldCompleted_WhenMultipleGoroutineDoNormalAddAndDone(t *testing.T) {
	wg := sync.WaitGroup{}

	loopCount := 100

	wg.Add(loopCount)

	for i := 0; i < loopCount; i++ {
		go func(i int) {
			println("executed ", i, "")
			wg.Done()
		}(i)
	}

	isSuccess := true

	wg.Wait()

	assert.True(t, isSuccess)

}

// TestWaitGroup_ShouldCompleted_WhenReuseWaitGroupAfterWaitReturnedAtFirstRound 用于测试在 WaitGroup 对象的 Wait 方法返回后，
// 是否可以再次复用 WaitGroup 对象来等待新的并发的 Goroutine。
// 首先，我们创建了 loopCount 个并发的 Goroutine，每个 Goroutine 完成任务之后都会调用 Done 方法。
// 在主 Goroutine 中我们调用 Wait 方法来阻塞等待这一组 Goroutine 的完成。
// 如果所有的 Goroutine 完成任务并调用了 Done 方法，那么 Wait 方法应该能够返回，说明第一轮的并发操作成功完成。
// 在第一轮的并发操作完成后，我们再次复用了这个 WaitGroup 对象，再次添加 loopCount 个并发的 Goroutine 并在主 Goroutine 中调用 Wait 方法阻塞等待。
// 如果这一轮的并发操作也能够成功完成，那么说明 WaitGroup 对象在 Wait 方法返回后可以被再次复用。
func TestWaitGroup_ShouldCompleted_WhenReuseWaitGroupAfterWaitReturnedAtFirstRound(t *testing.T) {
	wg := sync.WaitGroup{}

	loopCount := 100

	// 第一轮
	wg.Add(loopCount)

	for i := 0; i < loopCount; i++ {
		go func(i int) {
			println("executed ", i, "")
			wg.Done()
		}(i)
	}

	isFirstRoundSuccess := false
	wg.Wait()
	isFirstRoundSuccess = true

	// 第二轮
	wg.Add(loopCount)

	for i := 0; i < loopCount; i++ {
		go func(i int) {
			println("executed ", i, "")
			wg.Done()
		}(i)
	}

	isNextRoundSuccess := false

	wg.Wait()

	isNextRoundSuccess = true

	assert.True(t, isFirstRoundSuccess)
	assert.True(t, isNextRoundSuccess)
}

// TestWaitGroup_ShouldAllReturned_WhenMultipleGoroutineWaitForDone 测试了多个Goroutine
// 在同一时刻等待一个WaitGroup。在这个测试中，添加了一个Goroutine到WaitGroup并且该Goroutine会在
// 一段时间后调用Done()方法。而在其他的Goroutine中，它们通过调用Wait()方法等待这个WaitGroup。
// 所有等待该WaitGroup的Goroutine都应该在Done被调用后解锁。如果所有的Goroutine都顺利解锁并且Wait方法返回，
// 则表明WaitGroup可以正确地等待所有的Goroutine完成。
func TestWaitGroup_ShouldAllReturned_WhenMultipleGoroutineWaitForDone(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		wg.Done()
		println("Done")

	}()

	go func() {
		wg.Wait()
		println("Wait complete 1")
	}()

	go func() {
		wg.Wait()
		println("Wait complete 2")
	}()
	isSuccess := false
	wg.Wait()
	isSuccess = true

	assert.True(t, isSuccess)
}

// TestWaitGroup_ShouldBlocked_WhenAddCountGreaterThanDoneCount  测试了一种特殊情况
// 当我们添加了更多的需要等待的 goroutine 比实际完成的 goroutine 时，WaitGroup 的 Wait
// 方法应该阻塞，直至等待的 goroutine 完成。在这个测试用例中，我们专门添加了一个超时
// 机制来判断 Wait 是否正常阻塞。
func TestWaitGroup_ShouldBlocked_WhenAddCountGreaterThanDoneCount(t *testing.T) {
	wg := sync.WaitGroup{}
	ch := make(chan bool)

	wg.Add(2)
	go func() {
		println("Hello WaitGroup")
		wg.Done()
	}()

	go func() {
		wg.Wait()
		ch <- true
	}()

	isCompleted := false
	isTimeout := false
	select {
	case <-ch:
		isCompleted = true
	case <-time.After(10 * time.Second):
		isTimeout = true
	}

	assert.False(t, isCompleted)
	assert.True(t, isTimeout)
}

// TestWaitGroup_ShouldPanic_WhenAddNegtiveCounterWithoutAddPositiveCounter 在添加负数
// goroutine 时，WaitGroup 应该 panic。这是因为添加负数的 goroutine 没有任何意义，反而
// 会导致等待的 goroutine 数量变为负数，这是一种错误的使用方式。
func TestWaitGroup_ShouldPanic_WhenAddNegtiveCounterWithoutAddPositiveCounter(t *testing.T) {
	wg := sync.WaitGroup{}

	assert.Panics(t, func() { wg.Add(-1) })
}

// TestWaitGroup_ShouldPainic_WhenTryToReusePreviousWaitHasReturned 当我们尝试在上一个
// Wait 调用还没有返回时，复用 WaitGroup，WaitGroup 应该 panic。
// 这是因为如果我们在上一个Wait 调用没有返回之前就复用 WaitGroup，那么新增的 goroutine 可能会影响上一个 Wait 调用
// 的正确返回结果。
func TestWaitGroup_ShouldPainic_WhenTryToReusePreviousWaitHasReturned(t *testing.T) {
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(3)

	wg.Add(1)

	go func() {
		fmt.Println("Wait executed")
		wg.Wait()
		fmt.Println("Wait completed")
	}()

	go func() {
		time.Sleep(time.Second * 5)
		wg.Done()
		fmt.Println("Add executed")
		wg.Add(1)
	}()

	go func() {
		time.Sleep(time.Second * 5)
		fmt.Println("Done executed")
		wg.Done()
	}()

	time.Sleep(100 * time.Second)
}

// TestWaitGroup_ShouldPanic_WhenNoAddCalled 用于测试在没有通过 Add 方法添加任何需要
// 等待的 Goroutine 的情况下，直接调用 Done 方法是否会抛出 panic。
// 根据 WaitGroup 的使用规则，我们必须显式地通过调用 Add 方法告诉 WaitGroup，有多少个 Goroutine 需要等待完成。
// 如果我们没有通过 Add 方法添加任何需要等待的 Goroutine，而直接调用 Done 方法，那么应该
// 会抛出 panic，因为 Done 方法预期会有一些任务需要它去完成，然而这些任务并未被正确添加到
// WaitGroup 中。
func TestWaitGroup_ShouldPanic_WhenNoAddCalled(t *testing.T) {
	wg := sync.WaitGroup{}

	assert.Panics(t, func() {
		wg.Done()
	})

}
