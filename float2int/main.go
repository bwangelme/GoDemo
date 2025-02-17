package main

import (
	"fmt"
	"math"
)

// 将一个 float 类型的数字乘以 100，并且舍去其他小数位（即进行向下取整）
func main() {
	// 示例浮动数
	num := 3.14976

	// 将 num 乘以 100 并取整
	result := int(math.Floor(num * 100))

	// 输出结果
	fmt.Println(result) // 输出：314
}
