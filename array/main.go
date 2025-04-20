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

/*
测试修改 slice 底层的 array,会影响上层的 slice
*/
func arrSlice() {
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

/*
测试，子函数中改变数组，也会影响父函数中的数组
*/
func swapArr(arr []int64) {
	for idx := 0; idx < len(arr)/2; idx++ {
		tmp := arr[idx]
		arr[idx] = arr[len(arr)-1-idx]
		arr[len(arr)-1-idx] = tmp
	}
	fmt.Println("Swap sub arr", arr)
}

func arrSwap() {
	var arr = []int64{2, 3, 4, 5, 6, 7, 8}
	fmt.Println("Arr before swap", arr)
	swapArr(arr[2:])
	fmt.Println("Arr after swap", arr)
}

func main() {
	arrSlice()
}
