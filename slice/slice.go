package main

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	EmptySlice()
}

func EmptySlice() {
	// nil 切片的底层数组是 nil
	// 空切片的底层数组是一个 0 值数组
	// 所有空切片的底层都指向同一个数组
	var s1 []int
	var s2 = make([]int, 0)
	var s4 = make([]int, 0)

	fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)), *(*reflect.SliceHeader)(unsafe.Pointer(&s2)), *(*reflect.SliceHeader)(unsafe.Pointer(&s4)))
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s1))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data)
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s4))).Data)

	fmt.Printf("%v\n", unsafe.SliceData(s1))
	fmt.Printf("%v\n", unsafe.SliceData(s2))
	fmt.Printf("%v\n", unsafe.SliceData(s4))
}

func sliceOp() {
	// 切片的end 不能超过 slice 的长度
	i := []int{1, 2, 3, 4, 5}
	fmt.Println(i[2:])

	// 参考文章 https://mp.weixin.qq.com/s/fgy3FnfvFgRMjQ8in_otjA
	input := []byte("Hello, World")
	// 大多数匿名值都是不可寻址的(复合字面值是个例外)
	// 对数组进行切片操作要求该切片是可寻址的
	// 故以下代码会有编译错误
	//hash := sha1.Sum(input)[:5]

	// 将 Sum 的返回值赋值给一个临时变量后，它就是可寻址的变量了，可以使用切片操作
	v := sha1.Sum(input)
	hash := v[:5]
	fmt.Println(hash)
}
