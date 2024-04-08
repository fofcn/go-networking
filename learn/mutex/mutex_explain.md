# sync.Mutex

## 结构体
```go
type Mutex struct {
	state int32
	sema  uint32
}
```

## Lock()

加锁，如果锁没有被占用，则直接加锁，如果锁被占用，则阻塞等待。

先了解下go的CAS操作，CAS操作是原子操作，保证操作的原子性。
```go

// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
// Consider using the more ergonomic and less error-prone [Int32.CompareAndSwap] instead.
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)

```

go中的加锁的实现。
Lock由第一个Goroutine调用时可以使用CAS快速加锁返回。
其他Goroutine调用时都尝试使用CAS来加锁，因为此时之前的Goroutine可能已经解锁了
如果其他Goroutine没有解锁，那么就调用慢速加锁过程lockSlow。
```go

func (m *Mutex) Lock() {
	// 先尝试加锁，CAS操作；
    // 这里算是一个抢占的操作，如果锁没有被占用，则直接加锁，如果锁被占用，则阻塞等待。
    // 抢占操作可以提高性能，但是会增加等待时间。
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}

```

LockSlow()也会首先自旋来尝试获取锁，如果经过几次自旋失败后才会。。。
### lockSlow()
```go

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		// 如果有锁定并且可以自旋，那么就进入自旋；尝试N次后退出,我个人测试4次就退出了
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		
        // 检测Muxtex是否处于饥饿模式
        // old&mutexStarving == 0不处于饥饿模式
		if old&mutexStarving == 0 {
            // 将new变量设置为已锁定状态
			new |= mutexLocked
		}

        // 增加锁等待者计数
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}

        // 设置为饥饿模式
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		if awoke {
			// The goroutine has been woken from sleep,
			// so we need to reset the flag in either case.
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
			}
			// If we were already waiting before, queue at the front of the queue.
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			runtime_SemacquireMutex(&m.sema, queueLifo, 1)
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			if old&mutexStarving != 0 {
				// If this goroutine was woken and mutex is in starvation mode,
				// ownership was handed off to us but mutex is in somewhat
				// inconsistent state: mutexLocked is not set and we are still
				// accounted as waiter. Fix that.
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					// Exit starvation mode.
					// Critical to do it here and consider wait time.
					// Starvation mode is so inefficient, that two goroutines
					// can go lock-step infinitely once they switch mutex
					// to starvation mode.
					delta -= mutexStarving
				}
				atomic.AddInt32(&m.state, delta)
				break
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}

	if race.Enabled {
		race.Acquire(unsafe.Pointer(m))
	}
}

```