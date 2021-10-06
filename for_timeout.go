package main

import (
	"fmt"
	"reflect"
	"time"
)

func worker(i int, res chan []*int) {
	time.Sleep(200 * time.Millisecond)

	close(res)
	//res <- i
}

func main() {
	t := time.NewTimer(900 * time.Millisecond)
	cnt := 1
	var v []*int
	fmt.Println(v, v == nil, len(v))
	var result = make(chan []*int)
	for i := 0; i < cnt; i++ {
		go worker(i, result)
	}

LOOP:
	for i := 0; i < cnt; i++ {
		select {
		case <-t.C:
			fmt.Printf("%v timeout\n", i)
			break LOOP
		case v = <-result:
			fmt.Println(v, reflect.TypeOf(v), len(v), v == nil)
			fmt.Printf("%v return %v\n", i, v)
		}
	}
}
