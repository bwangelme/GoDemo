package main

/*
getName 函数想从 input 中拿到 . 分割的最后一组字符串

使用 `len(tokens)-1` 的索引方式，保证一定能够拿到内容，不会引发 panic
*/

import (
	"fmt"
	"strings"
)

func getName(input string) string {
	var tokens = strings.Split(input, ".")
	return tokens[len(tokens)-1]
}

func main() {
	fmt.Println(getName("abc"))     // == "abc"
	fmt.Println(getName("abc.edf")) // == "edf"
	fmt.Println(getName(""))        // == ""
}
