package main

import (
	"fmt"
	"time"
)

func timeoutWorker(c chan int) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		c <- i
	}
}

func main() {

	numChan := make(chan int, 3)

	go timeoutWorker(numChan)
	for i := 0; i < 10; {
		//fmt.Println("Check")
		select {
		case n := <-numChan:
			i++
			fmt.Println(n)
		default:

		}
	}
}
