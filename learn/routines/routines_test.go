package routines_test

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 获取Go协程使用的线程数量
func TestGPROC_ShouldReturnDefaultNumer_WhenNotSetProcNumber(t *testing.T) {
	procnum := runtime.GOMAXPROCS(0)
	fmt.Printf("default proc number: %d\n", procnum)
}

// runtime设置Go协程使用的线程数量
func TestGPROC_ShouldReturnSpecificNumer_WhenSetProcNumber(t *testing.T) {
	specnum := 4
	runtime.GOMAXPROCS(specnum)
	fmt.Printf("set proc number: %d\n", specnum)

	assert.Equal(t, specnum, runtime.GOMAXPROCS(0))
}

// 标准Go函数
func standardFunc(ch chan bool) {
	println("Hello, Standard Function Go Routine")
	ch <- true
}

// 标准函数创建协程
func TestRoutine_ShouldSuccess_WhenCreateWithStandardFunction(t *testing.T) {
	ch := make(chan bool)

	// func为标准函数
	go standardFunc(ch)

	ret := <-ch
	assert.True(t, ret)
}

// 闭包/匿名函数创建协程
func TestRoutine_ShouldSuccess_WhenCreateWithAnonymousFunction(t *testing.T) {
	ch := make(chan bool)

	// func为闭包/匿名函数
	go func() {
		println("Hello, Anonymous Function Go Routine")
		ch <- true
	}()

	ret := <-ch
	assert.True(t, ret)
}

type s struct {
	ch chan bool
}

func (s *s) run() {
	println("Hello, Struct Method Go Routine")
	s.ch <- true
}

func (s *s) wait() {
	<-s.ch
}
func TestRoutine_ShouldSuccess_whenCreateWithStructMethod(t *testing.T) {
	// 定义struct变量
	s := &s{
		ch: make(chan bool),
	}

	// 创建协程
	go s.run()

	// 等待执行完成
	s.wait()
}

func scheduleFunc(wg *sync.WaitGroup, f interface{}, args ...interface{}) {
	funcVal := reflect.ValueOf(f)

	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}

	// 创建新的 goroutine 时 WaitGroup 计数加 1
	wg.Add(1)

	go func() {
		defer wg.Done() // goroutine 执行结束后 WaitGroup 计数减 1
		funcVal.Call(in)
	}()
}

func task1(a string) {
	fmt.Printf("Hello: %s\n", a)
}

func task2(a, b string) {
	fmt.Printf("Hello: %s-%s\n", a, b)
}

func TestRoutine_ShouldSuccess_whenCreateWithReflect(t *testing.T) {
	var wg sync.WaitGroup // 创建一个 WaitGroup

	scheduleFunc(&wg, task1, "Hello, goroutine!")
	scheduleFunc(&wg, task2, "Hello", "goroutine!")

	wg.Wait() // 等待所有 goroutine 结束

}

// 不用太关注api和语法，只需要知道每个一秒钟打印"Hello background task"
func backgroundTask(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // 接收到取消信号，结束 goroutine
			return
		case <-ticker.C: // 每次 ticker 到时，打印一条消息
			println("Hello background task")
		}
	}
}
func TestRoutine_ShouldStop_whenSendCancelWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go backgroundTask(ctx)
	// 让 协程 运行一段时间
	time.Sleep(time.Second * 5)
	// 发送取消信号
	cancel()
	// 给协程留一点时间处理信号
	time.Sleep(time.Second * 2)
}

func signaltask(ch chan bool) {
	for {
		select {
		// 接收到取消信号，结束协程
		case <-ch:
			return
			// 没有接收到取消信号，打印一条消息
		default:
			println("Hello signal task")
			time.Sleep(time.Second * 1)
		}
	}
}
func TestRoutine_ShouldStop_WhenSendCancelSignal(t *testing.T) {
	ch := make(chan bool)
	go signaltask(ch)
	// 让协程运行5秒钟
	time.Sleep(time.Second * 5)
	// 发送取消信号
	ch <- true
	// 给协程留一点时间处理信号
	time.Sleep(time.Second * 2)
}

func TestRoutine_ShouldBeKilled_WhenTooManyChannelsAreCreated(t *testing.T) {
	for i := 0; i < math.MaxInt64; i++ {
		go func() {
			println("Hello, goroutine!")
			time.Sleep(1000 * time.Second)
		}()
	}
}
