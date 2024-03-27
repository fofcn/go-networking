package channel

import (
	"sync"
	"sync/atomic"
	"testing"

	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
)

type RateLimiter struct {
	tokens       chan struct{}
	refillTicker *time.Ticker
	closeCh      chan struct{}
}

func NewRateLimiter(rate int) *RateLimiter {
	r := &RateLimiter{
		tokens:       make(chan struct{}, rate),
		refillTicker: time.NewTicker(time.Second / time.Duration(rate)),
		closeCh:      make(chan struct{}),
	}

	go r.refill()

	return r
}

func (r *RateLimiter) refill() {
	for {
		select {
		case <-r.refillTicker.C:
			select {
			case r.tokens <- struct{}{}:
			default:
			}
		case <-r.closeCh:
			r.refillTicker.Stop()
			return
		}
	}
}

// Attempt to acquire a token, return false if there are none available
func (r *RateLimiter) Acquire() {
	<-r.tokens
}

// Attempt to acquire a token, return false if there are none available
func (r *RateLimiter) TryAcquire() bool {
	select {
	case <-r.tokens:
		return true
	default:
		return false
	}
}

// Close the RateLimiter and release all resources
func (r *RateLimiter) Close() {
	close(r.closeCh)
}

func myTask(id int) {
	fmt.Printf("time: %v workder %d is working\n", time.Now(), id)
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("time: %v workder %d has finished\n", time.Now(), id)
}

func TestRateLimiter_ShouldPermitWithBlocking_WhenRequestOnce(t *testing.T) {
	rateLimiter := NewRateLimiter(100)

	startTime := time.Now()
	for i := 0; i < 1; i++ {
		rateLimiter.TryAcquire()
		myTask(i)
	}
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %v\n", elapsedTime)
	fmt.Printf("explect time: %v\n", 300*time.Millisecond)
	assert.True(t, elapsedTime < 300*time.Millisecond)
}

func TestRateLimiter_ShouldLimitPermits_WhenGivenLimitedResource(t *testing.T) {
	var counter int32 = 0
	rateLimiter := NewRateLimiter(100)
	wg := sync.WaitGroup{}
	startTime := time.Now()
	for i := range 1000 {
		wg.Add(1)
		go func() {
			rateLimiter.Acquire()
			myTask(i)
			atomic.AddInt32(&counter, 1)
			wg.Done()
		}()

	}
	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %v\n", elapsedTime)
	fmt.Printf("should greater than explect time: %v\n", 10*time.Second)
	assert.Equal(t, counter, int32(1000))
	assert.True(t, 10*time.Second < elapsedTime)
}
