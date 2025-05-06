package main

import (
	"fmt"
	"sync"
)

/*
p 1-100 3个
c print

## 程序思路

1. 一共有3个 channel，每个生产者 P 选择一个 channel 进行写入
2. 生产者 P 写完以后，关闭 channel
3. 消费者 C 从 三个 channel 中选择一个进行消费, 如果当前 channel 已经关闭了，那么设置为 nil, 让 select 不再读取
3.1 消费者 C 判断三个 channel 都关闭了，协程退出
4. main 等待生产者和消费者都结束后，程序结束
*/

func P(wg *sync.WaitGroup, name string, ch chan string) {
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("[%v]: %d", name, i)
		ch <- msg
	}
	close(ch)
	wg.Done()
}

func C(wg *sync.WaitGroup, chs []chan string) {
	defer wg.Done()
Out:
	for {
		select {
		case msg1, ok := <-chs[0]:
			if !ok {
				chs[0] = nil
			} else {
				fmt.Println(msg1)
			}
		case msg2, ok := <-chs[1]:
			if !ok {
				chs[1] = nil
			} else {
				fmt.Println(msg2)
			}
		case msg3, ok := <-chs[2]:
			if !ok {
				chs[2] = nil
			} else {
				fmt.Println(msg3)
			}
		}

		if chs[0] == nil && chs[1] == nil && chs[2] == nil {
			break Out
		}
	}
	fmt.Println("Consumer stop")
}

func main() {
	chs := make([]chan string, 0)
	for i := 0; i < 3; i++ {
		ch := make(chan string)
		chs = append(chs, ch)
	}
	var wg = sync.WaitGroup{}
	wg.Add(5)

	go P(&wg, "w1", chs[0])
	go P(&wg, "w2", chs[1])
	go P(&wg, "w3", chs[2])

	go C(&wg, chs)
	go C(&wg, chs)

	wg.Wait()
}
