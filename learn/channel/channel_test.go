package channel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go channel 是go routines通信的桥梁
// Go Channel的基本使用：
//  创建：
// ch := make(chan int)
// 发送数据到Channel
// ch <- 1
// 从Channel接受数据
// val <- ch
// <- 符号
//
// go 中chan分为无缓冲Channel和有缓冲Channel
// 无缓冲Channel
//    无缓冲Channel没有存储数据的能力
//    发送方向Channel中发送数据的时候，发送方会阻塞直到有接受者接受这个数据
//    无缓冲Channel典型应用就是go协程同步通信
//    无缓冲Channel保证通信双方都要准备好数据交换
// 有缓冲Channel
//    有缓冲Channel需要定义Channel的容量
//    发送方向有缓冲Channel发送数据的时候，只有容量满的时候才会阻塞
//    接收方只有在有缓冲Channel为空时才会阻塞
//    有缓冲通道的典型应用场景是生产者和消费者

func TestUnbufferedChannel_ShouldPanic_whenWriteValueToAClosedChannel(t *testing.T) {

	f := func() {
		ch := make(chan int)
		close(ch)

		ch <- 1
	}

	assert.Panics(t, f, "should panic")
}

func TestUnbufferedChannel_ShouldSuccess_whenRecvValueAtAClosedChannel(t *testing.T) {
	ch := make(chan int)
	close(ch)
	var val = <-ch
	assert.Equal(t, 0, val)
}

func TestUnbufferedChannel_ShouldSuccess_whenRecvEmptyValueFromAClosedChannel(t *testing.T) {
	ch := make(chan string)
	close(ch)
	var val = <-ch
	assert.Equal(t, "", val)
	val = <-ch
	assert.Equal(t, "", val)
}

func TestUnbufferedChannel_ShouldReturnNil_whenRecvDataFromAClosedChannel(t *testing.T) {
	var expectedStr *string = nil
	ch := make(chan *string)
	close(ch)
	var val = <-ch
	assert.Equal(t, expectedStr, val)
	val = <-ch
	assert.Equal(t, expectedStr, val)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	// 将累加结果发送到channel
	c <- sum
}

func TestUnbufferedChannel_ShouldRecvValues_WhenWriteValueToChannel(t *testing.T) {

	// 创建无缓冲channel
	c := make(chan int)

	// given
	s := []int{1, 2, 3, 4, 5, 6}

	// when
	// 执行数组累加
	go sum(s[:], c)
	ret1 := <-c

	// should
	// 和应该是21
	assert.Equal(t, 21, ret1)
}

func TestUnbufferedChannel_ShouldReadTimeout_WhenNoValueWriteToChannel(t *testing.T) {

	// given
	c := make(chan int)
	is_timeout := false

	// when
	select {
	case <-c:
	case <-time.After(3 * time.Second):
		// should
		is_timeout = true
	}

	assert.True(t, is_timeout)

}

func TestUnbufferedChannel_ShouldWriteTimeout_WhenNoRoutineReadTheChannel(t *testing.T) {

	// given
	c := make(chan int)
	is_timeout := false
	try_to_write_value := 1

	// when
	select {
	case c <- try_to_write_value:
	case <-time.After(3 * time.Second):
		// should
		is_timeout = true
	}

	assert.True(t, is_timeout)

}

func TestBufferedChannel_ShouldNotBlock_WhenWriteValueLessThanCapacity(t *testing.T) {
	// given
	c := make(chan int, 2)

	// when
	c <- 1
	c <- 2

	// then
	assert.Equal(t, 1, <-c)
	assert.Equal(t, 2, <-c)
}

func TestBufferedChannel_ShouldBlock_WhenCapacityReached(t *testing.T) {
	// given
	c := make(chan int, 2)

	// when
	c <- 1
	c <- 2

	// this will block channel
	is_block := false
	go func() {
		defer func() {
			if r := recover(); r != nil {
				is_block = true
			}
		}()
		c <- 3
	}()

	time.Sleep(1 * time.Second) // wait for goroutine completes
	assert.True(t, is_block)
}

// 先定义一个worker函数
// worker函数从无缓冲channel中接收
// 可以接到到数据就执行后面的打印内容
// 打印完成后退出
func worker(id int, lock chan bool) {
	var shouldRun = <-lock
	if shouldRun {
		fmt.Printf("time: %v Worker %d is working\n", time.Now(), id)
		time.Sleep(time.Second)
		fmt.Printf("time: %v Worker %d has finished\n", time.Now(), id)
	}
}

func TestUnbufferedChannel_ShouldRunOneByOne_When(t *testing.T) {
	lock := make(chan bool, 1)

	// 启动5个goroutine等待释放接收
	for i := 0; i < 5; i++ {
		go worker(i, lock)
	}

	// 发送5个true到channel
	for i := 0; i < 5; i++ {
		lock <- true
		time.Sleep(time.Second)
	}

	close(lock)

	time.Sleep(10 * time.Second)
}

func future(id int, delay time.Duration, resChan chan int) {
	time.Sleep(delay)
	fmt.Printf("Hi, I have finished my task, my id is %d\n", id)
	resChan <- id
}

func anyOf(futures ...<-chan int) <-chan int {
	result := make(chan int)
	for _, future := range futures {
		go func(f <-chan int) {
			result <- <-f
		}(future)
	}
	return result
}

func TestAnyOf_ShouldSuccess(t *testing.T) {
	// 创建无缓冲的 channel
	resChan1 := make(chan int)
	resChan2 := make(chan int)
	resChan3 := make(chan int)

	// 启动 goroutines
	go future(1, 3*time.Second, resChan1)
	go future(2, 2*time.Second, resChan2)
	go future(3, 5*time.Second, resChan3)

	result := anyOf(resChan1, resChan2, resChan3)

	assert.Equal(t, 2, <-result)
}
