package main

import "fmt"

func reverse(arr []int64) []int64 {
	for idx := 0; idx < len(arr)/2; idx++ {
		tmp := arr[idx]
		arr[idx] = arr[len(arr)-1-idx]
		arr[len(arr)-1-idx] = tmp
	}

	return arr
}

func main() {
	arr := []int64{3, 5, 8, 9, 4}
	arr1 := reverse(arr)

	arr[0] = -1

	fmt.Printf("arr: %v\narr1: %v\n", arr, arr1)
}
