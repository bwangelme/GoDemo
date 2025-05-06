package main

import (
	"fmt"
	"sync"
)

/*
p 1-100 3个
c print
*/

func P(wg *sync.WaitGroup, name string, ch chan string) {
	for i := 0; i < 2; i++ {
		msg := fmt.Sprintf("[%v]: %d", name, i)
		ch <- msg
	}
	close(ch)
	wg.Done()
}

func C(wg *sync.WaitGroup, chs []chan string) {
Out:
	for {
		select {
		case msg1, ok := <-chs[0]:
			fmt.Printf(msg1)
			if !ok {
				break Out
			}
		case msg2, ok := <-chs[1]:
			fmt.Printf(msg2)
			if !ok {
				break Out
			}
		case msg3, ok := <-chs[2]:
			fmt.Printf(msg3)
			if !ok {
				break Out
			}
		}
	}
	wg.Done()
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
