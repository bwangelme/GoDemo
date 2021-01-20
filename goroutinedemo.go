package main

import (
	"context"
	"fmt"
	"time"
)

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	close(t.stop)
}

// Shutdown 函数用来关闭运行着的 Goroutine，并且还可以设置关闭的超时时间
func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
		// 关闭 Goroutine 成功
	case <-ctx.Done():
		// 关闭 Goroutine 超时
	}
}

// 这个例子就体现了使用 Goroutine 的原则:
// 1. 调用者在创建了 Goroutine 后，一定能够控制它什么时候退出
// 2. 知道 Goroutine 什么时候退出
// 3. 一定要处理 Goroutine 的超时错误
// 4. 将并发逻辑交给调用者，即让调用者来创建 Goroutine
func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "foo")
	_ = tr.Event(context.Background(), "bar")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tr.Shutdown(ctx)
}
