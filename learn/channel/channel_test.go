package channel_test

import (
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

func TestUnbufferedChannel_ShouldRecvValues_WhenWriteValueToChannel(t *testing.T) {

	c := make(chan int)

	// given
	s := []int{1, 2, 3, 4, 5, 6}

	// when
	go sum(s[:], c)
	ret1 := <-c

	// should
	assert.Equal(t, 21, ret1)
}

func TestUnbufferedChannel_ShouldTimeout_WhenNoValueWriteToChannel(t *testing.T) {

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

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}
