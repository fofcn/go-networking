package network_test

import (
	"go-networking/network"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewResponseFutureShouldSetCorrectTimestampWhenInitialized 测试 NewResponseFuture 是否在初始化时设置了正确的时间戳
func TestNewResponseFutureShouldSetCorrectTimestampWhenInitialized(t *testing.T) {
	seq := uint64(123)
	rf := network.NewResponsePromise(seq, 30*time.Second)
	assert.WithinDuration(t, time.Now(), rf.Timestamp(), time.Millisecond*500, "Timestamp should be within 500 milliseconds of now")
	rf.Close() // 清理资源
}

// TestResponseFutureShouldReturnFrameWhenAddIsCalled 测试 Add 是否正确返回 Frame 对象
func TestResponseFutureShouldReturnFrameWhenAddIsCalled(t *testing.T) {
	seq := uint64(123)
	frame := &network.Frame{Seq: seq}
	rf := network.NewResponsePromise(seq, 30*time.Second)

	go func() {
		rf.Add(frame)
	}()

	resultFrame, err := rf.Wait()
	assert.NoError(t, err, "Wait should not return an error")
	assert.Equal(t, frame, resultFrame, "The frame returned by Wait should be the same as the one added")
	rf.Close() // 清理资源
}

// TestResponseFutureShouldReturnErrorWhenWaitTimeout 测试 Wait 是否在超时的情况下返回错误
func TestResponseFutureShouldReturnErrorWhenWaitTimeout(t *testing.T) {
	seq := uint64(123)
	rf := network.NewResponsePromise(seq, 50*time.Millisecond) // 使用较短的超时时间以便测试

	resultFrame, err := rf.Wait()
	assert.Nil(t, resultFrame, "The result frame should be nil on timeout")
	assert.Error(t, err, "An error should be returned when waiting times out")
	assert.Contains(t, err.Error(), "waiting for response timeout", "The error message should contain the timeout information")
	rf.Close() // 清理资源
}

// TestResponseFutureShouldSupportMultipleConcurrentWaits 测试 Wait 是否可以支持多个并发等待
func TestResponseFutureShouldSupportMultipleConcurrentWaits(t *testing.T) {
	seq := uint64(123)
	frame := &network.Frame{Seq: seq}
	rf := network.NewResponsePromise(seq, 5*time.Second)

	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	var resultFrame1, resultFrame2 *network.Frame

	// 第一个并发 Wait
	go func() {
		defer wg.Done()
		resultFrame1, err1 = rf.Wait()
	}()

	// 第二个并发 Wait
	go func() {
		defer wg.Done()
		resultFrame2, err2 = rf.Wait()
	}()

	// 给 Wait 足够的时间启动
	time.Sleep(100 * time.Millisecond)
	rf.Add(frame) // 解除 Wait 阻塞

	wg.Wait()
	assert.NoError(t, err1, "First concurren.Wait() should not return an error")
	assert.Equal(t, frame, resultFrame1, "The frame returned by the first concurrent Wait should be the same as the one added")
	assert.NoError(t, err2, "Second concurrent Wait() should not return an error")
	assert.Equal(t, frame, resultFrame2, "The frame returned by the second concurrent Wait should be the same as the one added")
	rf.Close() // 清理资源
}
