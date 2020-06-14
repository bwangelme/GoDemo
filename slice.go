package main

import "fmt"

func main() {
	// 切片的end 不能超过 slice 的长度
	i := []int{1, 2, 3, 4, 5}
	fmt.Println(i[2:])
}
