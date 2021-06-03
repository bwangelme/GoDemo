package main

import (
	"fmt"
	"time"
)

func main() {
	timer := time.NewTimer(5 * time.Second)
	c := make(chan int)

	go func() {
		select {
		case c <- 1:
			fmt.Println("Write")
		default:
			fmt.Println("No One")
		}
	}()

	for {
		select {
		case <-timer.C:
			fmt.Println("End")
			return
		case val := <-c:
			fmt.Println("Read Val", val)
		}
	}
}
