package main

// 这段代码展示了 Context 的使用场景
// 主进程可以直接停止所有的子进程

import (
	"context"
	"fmt"
	"time"
)

func Monitor(ctx context.Context, num int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%d号监控程序退出.\n", num)
			return
		default:
			fmt.Printf("%d 正在监控中...\n", num)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 3; i++ {
		go Monitor(ctx, i)
	}

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(4 * time.Second)
	fmt.Printf("主进程退出\n")
}
