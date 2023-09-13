package main

import (
	"context"
	"fmt"
	"time"
)

/**
本程序测试如果存在多个 timeout, 那个先生效

deadlien time 在前面的那个生效
*/

func Monitor(ctx context.Context, num int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%d号监控程序退出.\n", num)
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Printf("%d 正在监控中...\n", num)
		}
	}
}
func main() {
	ctx1, cancel1 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel1()
	ctx2, cancel2 := context.WithTimeout(ctx1, 4*time.Second)
	defer cancel2()

	for i := 1; i <= 3; i++ {
		go Monitor(ctx2, i)
	}

	time.Sleep(5 * time.Second)
	fmt.Printf("主进程退出\n")
}
