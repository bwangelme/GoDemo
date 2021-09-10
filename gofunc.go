package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

func Go(handler func()) {
	go func() {
		defer func() {
			if res := recover(); res != nil {
				stack := debug.Stack()
				fmt.Println(res)
				fmt.Println(string(stack))
			}
		}()

		handler()
	}()
}

func main() {
	for i := 0; i < 10; i++ {
		var a = i
		Go(func() {
			if a == 3 {
				log.Panicln(a)
			}
			fmt.Println("正常结束", a)
		})
	}

	time.Sleep(1 * time.Second)
}
