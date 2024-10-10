package main

import (
	"fmt"
	"time"
)

func main() {
	// 定义时间字符串
	timeStr := "20240829113745"

	// 解析时间字符串为 time.Time 类型
	layout := "20060102150405"
	t, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		fmt.Println("解析时间字符串时出错:", err)
		return
	}
	// 将时间转换为 Unix 时间戳
	timestamp := t.Unix()

	// 输出 Unix 时间戳
	fmt.Println("Unix 时间戳:", timestamp, t)
	fmt.Println("time now: ", time.Now().Unix())
}
