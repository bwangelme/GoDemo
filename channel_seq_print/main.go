package main

/*
## 程序说明

本文件中实现了一道 Go 面试题的解答，题目是：

创建 m 个 goroutine，他们依次打印 1-n 的数字，并且每个 goroutine 打印的数字都不同

例如，m = 3，n = 10，那么 3 个 goroutine 打印的数字分别是：
1, 4, 7, 10
2, 5, 8
3, 6, 9

每个 goroutine 打印的数字都不同，并且每个 goroutine 打印的数字都连续
*/

import (
	"fmt"
	"sync"
)

func worker(wg *sync.WaitGroup, idx int, channels []chan int, total int) {
	defer wg.Done()
	ch := channels[idx]
	nextCh := channels[(idx+1)%len(channels)]
	for v := range ch {
		fmt.Printf("Worker %d: %d\n", idx, v)
		if v >= total {
			break
		}
		nextCh <- v + 1
	}
	close(nextCh)
	fmt.Printf("Worker %d: exit\n", idx)
}

func main() {
	var (
		total = 10
		cnt   = 5
		wg    = sync.WaitGroup{}
	)

	channels := make([]chan int, cnt)
	for i := range cnt {
		channels[i] = make(chan int)
	}

	wg.Add(cnt)
	for i := range cnt {
		go worker(&wg, i, channels, total)
	}

	channels[0] <- 1
	wg.Wait()
}
