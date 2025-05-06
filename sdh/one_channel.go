package main

import (
	"fmt"
	"sync"
)

/*
## 程序思路

1. 只有一个 channel, 3个生产者和2个消费者都读写同一个  channel
2. 消费者根据 channel 是否 close，来决定协程结束，因此 main 中需要等生产者运行结束后，主动关闭 channel
3. 因此 wg 分成了 pwg 和 cwg 两个，pwg 判断生产者是否结束，cwg 判断消费者是否结束
*/

func Producer(wg *sync.WaitGroup, name string, ch chan string) {
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("[%v]: %d", name, i)
		ch <- msg
	}
	wg.Done()
}

func Consumer(wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println(v)
		}
	}
}

func main() {
	ch := make(chan string)
	var pwg = sync.WaitGroup{}
	var cwg = sync.WaitGroup{}

	pwg.Add(3)

	go Producer(&pwg, "w1", ch)
	go Producer(&pwg, "w2", ch)
	go Producer(&pwg, "w3", ch)

	go Consumer(&cwg, ch)
	go Consumer(&cwg, ch)

	// 生产者结束后，关闭 channel, 主动通知消费者退出
	pwg.Wait()
	close(ch)

	cwg.Wait()
}
