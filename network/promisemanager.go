package network

import (
	"go-networking/log"
	"sync"
	"time"
)

type PromiseM struct {
	rpTable map[uint64]ResponsePromise
	mux     sync.Mutex
}

func NewPromiseM() *PromiseM {
	return &PromiseM{
		rpTable: make(map[uint64]ResponsePromise),
	}
}

func (p *PromiseM) AddResp(frame *Frame) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if _, exists := p.rpTable[frame.Seq]; exists {
		// 尝试添加frame到promise中，这里假设Add方法可以抛出异常
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Recovered in AddResp: %v", err)
			}
		}()
		rp := p.rpTable[frame.Seq]
		rp.Add(frame)
	} else {
		log.Infof("what's wrong? frame sequence not matched with sequence no.: %d", frame.Seq)
	}
}

func (p *PromiseM) AddSeqPromise(seq uint64, rp ResponsePromise) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.rpTable[seq] = rp
}

func (p *PromiseM) DelSeqPromise(seq uint64) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if future, exists := p.rpTable[seq]; exists {
		// 关闭promise，这里假设Close方法可以抛出异常
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Recovered in DelSeqPromise: %v", err)
			}
		}()
		future.Close()
		delete(p.rpTable, seq)
	}
}

func (p *PromiseM) CloseRespPromis() {
	p.mux.Lock()
	defer p.mux.Unlock()

	for seq, future := range p.rpTable {
		delete(p.rpTable, seq)
		// 关闭CountDownLatch，这里假设Close方法可以抛出异常
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Recovered in CloseRespPromis: %v", err)
			}
		}()
		future.Close()
	}
}

func (p *PromiseM) CleanupRespPromise() {
	now := time.Now()
	p.mux.Lock()
	defer p.mux.Unlock()

	for seq, future := range p.rpTable {
		if now.Sub(future.Timestamp()) > 30*time.Second {
			// 如果ResponseFuture超过30秒钟
			// 从respTable删除，这里假设Close方法可以抛出异常
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("Recovered in CleanupRespPromise: %v", err)
				}
			}()
			delete(p.rpTable, seq)
			future.Close()
		}
	}
}
