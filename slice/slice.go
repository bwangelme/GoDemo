package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
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
