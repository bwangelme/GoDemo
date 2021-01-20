// single machine word data race
package main

import "fmt"

// 这个程序使用 go build -race racedemo/swracedemo_v2.go 编译后运行，就不会报错
// 因为 Ben 和 Jerry 的内存布局是一样的，但仍然可能存在数据v竞争的情况
// 即 interface 的 type 指向了 Ben，但实际的数据指向的是 Jerry

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("%v Says, \"Hello, my name is %v\"\n", b.name, b.name)
}

type Jerry struct {
	field1 *[5]byte
	field2 int
}

func (j *Jerry) Hello() {
	var bs = *j.field1
	var name = string(bs[:])
	fmt.Printf("%v Says, \"Hello, my name is %v\"\n", name, name)
}

func main() {
	var ben = &Ben{name: "Ben"}
	var jerry = &Jerry{field1: &[5]byte{'J', 'e', 'r', 'r', 'y'}, field2: 5}
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
