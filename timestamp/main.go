package main

import (
	"fmt"
	"time"
)

func main() {
	// 定义时间字符串
	timeStr := "2024-12-11 11:37:45.123"

	// 解析时间字符串为 time.Time 类型
	layout := "2006-01-02 15:04:05.999"
	t, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		fmt.Println("解析时间字符串时出错:", err)
		return
	}
	// 将时间转换为 Unix 时间戳
	timestamp := t.Unix()

	// 输出 Unix 时间戳
	fmt.Println("Unix 时间戳:", timestamp, t)
	fmt.Println("time now 时间戳: ", time.Now().Unix())
}
