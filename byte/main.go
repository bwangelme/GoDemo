package main

import "fmt"

/*
如果将一个不是标准字符串的 byte 转换成 string，会是什么输出

直接不输出内容
*/
func main() {
	data := []byte{
		0x00, 0x01, 0x03, 0x04,
	}
	fmt.Println("data str is", string(data))
}
