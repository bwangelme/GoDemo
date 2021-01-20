package main

import (
	"fmt"
	"sync"
	"time"
)

// 在编译程序时加入 race 选项，那么在程序运行时，如果检测到数据竞争，程序就会 panic
//ø> go build -race racedemo.go                                                                                                                                                                                      07:14:49 (01-17)
//ø> ./racedemo                                                                                                                                                                                                      07:14:58 (01-17)
//==================
//WARNING: DATA RACE
//Write at 0x00000064afc0 by goroutine 7:
//  main.Routine()
//      /home/xuyundong/Github/Golang/GoDemo/racedemo.go:27 +0x74
//
//Previous read at 0x00000064afc0 by goroutine 8:
//  main.Routine()
//      /home/xuyundong/Github/Golang/GoDemo/racedemo.go:24 +0x47
//
//Goroutine 7 (running) created at:
//  main.main()
//      /home/xuyundong/Github/Golang/GoDemo/racedemo.go:15 +0x75
//
//Goroutine 8 (finished) created at:
//  main.main()
//      /home/xuyundong/Github/Golang/GoDemo/racedemo.go:15 +0x75
//==================
//Final Counter: 2
//Found 1 data race(s)

var wg sync.WaitGroup
var Counter int = 0

func main() {
	for routine := 1; routine <= 2; routine++ {
		wg.Add(1)
		go Routine(routine)
	}

	wg.Wait()
	fmt.Printf("Final Counter: %d\n", Counter)
}

func Routine(id int) {
	for count := 0; count < 2; count++ {
		value := Counter
		time.Sleep(1 * time.Nanosecond)
		value++
		Counter = value
	}
	wg.Done()
}
