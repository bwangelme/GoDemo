package main

import (
	"fmt"
	"time"
)

func main() {
	// 定义时间字符串
	timeStr := "2023-10-05 14:30:45"

	// 定义布局字符串
	layout := "2006-01-02 15:04:05"

	// 解析时间字符串
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("解析时间出错:", err)
		return
	}

	// 获取当前时间
	now := time.Now()

	// 计算时间差
	diff := now.Sub(parsedTime)

	// 判断时间差是否不超过 5 分钟
	if diff.Abs() <= 5*time.Minute {
		fmt.Println("时间差不超过 5 分钟")
	} else {
		fmt.Println("时间差超过 5 分钟")
	}

	// 输出解析后的时间和当前时间
	fmt.Println("解析后的时间:", parsedTime)
	fmt.Println("当前时间:", now)
	fmt.Println("时间差:", diff)
}
