package main

import (
	"fmt"
)

/*
参数传递的时候，大都传递的是 slice, 传递数组需要写长度，太麻烦

修改切片的值会影响底层数组的值

往切片中新增元素，切片会生成新的数组，不会影响原来的底层数组
*/

type Item struct {
	Val int
}

func (s *Item) String() string {
	return fmt.Sprintf("<Item %v>", s.Val)
}

func main() {
	arr := make([]*Item, 10)
	arrHandler(arr)
	fmt.Println(arr)
	sliceHandler(arr)
	fmt.Println(arr)
}

func arrHandler(arr []*Item) {
	arr[3] = &Item{
		Val: 42,
	}
}

func sliceHandler(arr []*Item) {
	arr = append(arr, &Item{
		Val: 33,
	})
}
