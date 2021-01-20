// single machine word data race
package main

import "fmt"

// 这个程序可能会报 data race 错误
// maker 并不是一个 single machine word，interface 是由 type 和 data 两个指针组成的
// 所以输出的时候，可能会出现 在 maker 上出现 data race 的情况
// 出现下面这种错误，是因为 Jerry 和 Ben 的内存布局不一样

//ø> go build -race swracedemo.go                                                                                                                                                                                    07:39:40 (01-17)
//ø> ./swracedemo                                                                                                                                                                                                    07:38:12 (01-17)
//Ben Says, "Hello, my name is Ben"
//Ben Says, "Hello, my name is Ben"
//Ben Says, "Hello, my name is Ben"
//Jerry Says, "Hello, my name is Jerry"
//Jerry Says, "Hello, my name is Jerry"
//Jerry Says, "Hello, my name is Jerry"
//Jerry Says, "Hello, my name is Jerry"
//panic: runtime error: invalid memory address or nil pointer dereference
//[signal SIGSEGV: segmentation violation code=0x1 addr=0x14 pc=0x497d13]
//
//goroutine 1 [running]:
//fmt.(*buffer).writeString(...)
//        /home/xuyundong/.local/go/src/fmt/print.go:82
//fmt.(*fmt).padString(0xc0003e4520, 0x14, 0x541118)
//        /home/xuyundong/.local/go/src/fmt/format.go:110 +0xf9
//fmt.(*fmt).fmtS(0xc0003e4520, 0x14, 0x541118)
//        /home/xuyundong/.local/go/src/fmt/format.go:359 +0x6f
//fmt.(*pp).fmtString(0xc0003e44e0, 0x14, 0x541118, 0xc000000076)
//        /home/xuyundong/.local/go/src/fmt/print.go:447 +0x1c6
//fmt.(*pp).printArg(0xc0003e44e0, 0x522260, 0xc0003dde00, 0x76)
//        /home/xuyundong/.local/go/src/fmt/print.go:698 +0xd8f
//fmt.(*pp).doPrintf(0xc0003e44e0, 0x545f55, 0x20, 0xc00019fed0, 0x2, 0x2)
//        /home/xuyundong/.local/go/src/fmt/print.go:1030 +0x327
//fmt.Fprintf(0x566320, 0xc0001be008, 0x545f55, 0x20, 0xc00019fed0, 0x2, 0x2, 0x26, 0x0, 0x0)
//        /home/xuyundong/.local/go/src/fmt/print.go:204 +0x85
//fmt.Printf(...)
//        /home/xuyundong/.local/go/src/fmt/print.go:213
//main.(*Jerry).Hello(0xc0001b0020)
//        /home/xuyundong/Github/Golang/GoDemo/swracedemo.go:64 +0x134
//main.main()
//        /home/xuyundong/Github/Golang/GoDemo/swracedemo.go:87 +0x355

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	id   int
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("%v Says, \"Hello, my name is %v\"\n", b.name, b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("%v Says, \"Hello, my name is %v\"\n", j.name, j.name)
}

func main() {
	var ben = &Ben{name: "Ben", id: 20}
	var jerry = &Jerry{name: "Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}
