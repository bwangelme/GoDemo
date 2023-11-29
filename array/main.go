package main

import (
	"fmt"
	"reflect"
)

type Item struct {
	Val int
}

func (s *Item) String() string {
	return fmt.Sprintf("<Item %v>", s.Val)
}

func main() {
	var arr [10]*Item
	// 这里已经将数组转换成切片了
	arrHandler(arr[:])
	fmt.Println(reflect.TypeOf(arr))
	fmt.Println(arr)
}

func arrHandler(arr []*Item) {
	fmt.Println(reflect.TypeOf(arr))
	arr[3] = &Item{
		Val: 42,
	}
}
