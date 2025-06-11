package main

import "fmt"

// convertInt2Uint
// golang 在进行数字类型转换的时候，不会修改符号位，将数字的二进制位直接复制到目标类型的变量中
func convertInt2Uint() {
	var i1 int8 = -1
	var u1 = uint8(i1)
	fmt.Println("int8(-1) => uint8(255)")
	fmt.Println(u1)

	var u2 uint8 = 254
	var i2 = int8(u2)
	fmt.Println("uint8(254) => int8(-2)")
	fmt.Println(i2)
}

func main() {
	convertInt2Uint()
}
