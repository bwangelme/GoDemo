package main

import (
	"fmt"
	"time"
)

func main() {
	// 返回当前时间时间戳，时间戳是整数表示的秒
	t := time.Now().Unix()
	fmt.Println(t)
}
