package main

import (
	"fmt"
	"time"
)

/**
  测试 duration 可以和 0 比较，当 duration 是负数时，duration < 0 成立
 */
func main() {
	now := time.Now()
	future := now.Add(1 * time.Second)

	d := now.Sub(future)
	lessZero := d < 0

	fmt.Println(lessZero, d, future, now)
}
