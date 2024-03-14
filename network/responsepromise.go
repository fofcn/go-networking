package network

import (
	"fmt"
	"time"

	"github.com/quintans/toolkit/latch"
)

type ResponsePromise interface {
	Add(frame *Frame)
	Wait() (*Frame, error)
	Close()
	Timestamp() time.Time
}

type ResponsePromiseI struct {
	seq        uint64
	frame      *Frame
	createTime time.Time
	countdown  latch.CountDownLatch
}

func NewResponseFuture(seq uint64, timeout time.Duration) *ResponsePromiseI {
	rf := &ResponsePromiseI{
		countdown:  *latch.NewCountDownLatch(),
		createTime: time.Now(),
	}
	rf.countdown.Add(1)
	return rf
}

func (rf *ResponsePromiseI) Add(frame *Frame) {
	rf.frame = frame
	rf.countdown.Done()
}

func (rf *ResponsePromiseI) Wait() (*Frame, error) {
	isTimeout := rf.countdown.WaitWithTimeout(30 * time.Second)
	if isTimeout {
		return nil, fmt.Errorf("waiting for response timeout, seq: %d", rf.seq)
	}
	return rf.frame, nil
}

func (rf *ResponsePromiseI) Close() {
	rf.countdown.Close()
}

func (rf *ResponsePromiseI) Timestamp() time.Time {
	return rf.createTime
}
