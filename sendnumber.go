package main

/*
  这个程序中，两种写法向 ch 中写入的值不同，是因为 ch <- (运算式)，运算式求值的时机不同

  ch <- *num，运算式在主线程执行 <-ch 时才求值，所以求得的值是 789
  ch <- *num + 0 运算式在子线程运行的时候就求值了，求得的值是 123，等主线程执行 <-ch 将值写入到ch中
*/

import (
	"fmt"
	"time"
)

func sub(num *int, ch chan<- int) {
	//ch <- *num  // 写的值是789
	ch <- *num + 0 // 写入的值是123
}

func main() {
	var num = 123
	ch := make(chan int)
	go sub(&num, ch)

	time.Sleep(1 * time.Second)
	num = 789
	fmt.Println(<-ch)
}
