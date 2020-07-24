package main

import "fmt"

func main() {
	var nums []int = nil

	// for 循环可以迭代 nil 类型的数组容器，不会报错
	for _, i := range nums {
		fmt.Println(i)
	}
}
