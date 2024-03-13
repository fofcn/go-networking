package network

import (
	"fmt"
	"time"

	"github.com/quintans/toolkit/latch"
)

type ResponseFuture interface {
	Add(frame *Frame)
	Wait() (*Frame, error)
	Close()
}

type ResponseFutureI struct {
	seq       uint64
	frame     *Frame
	countdown latch.CountDownLatch
}

func NewResponseFuture(seq uint64, timeout time.Duration) *ResponseFutureI {
	rf := &ResponseFutureI{
		countdown: *latch.NewCountDownLatch(),
	}
	rf.countdown.Add(1)
	return rf
}

func (rf *ResponseFutureI) Add(frame *Frame) {
	rf.frame = frame
	rf.countdown.Done()
}

func (rf *ResponseFutureI) Wait() (*Frame, error) {
	isTimeout := rf.countdown.WaitWithTimeout(30 * time.Second)
	if isTimeout {
		return nil, fmt.Errorf("waiting for response timeout, seq: %d", rf.seq)
	}
	return rf.frame, nil
}

func (rf *ResponseFutureI) Close() {
	rf.countdown.Close()
}
