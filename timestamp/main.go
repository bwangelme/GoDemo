package main

import (
	"fmt"
	"time"
)

func Time2TS() {
	// 定义时间字符串
	timeStr := "2024-12-11 11:12:45.123"

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

func TS2Time() {
	ts := 1733887537
	t := time.Unix(1733887537, 0)
	fmt.Println("Unix 时间戳", ts)
	fmt.Println("时间戳转换后的本地时间", t)
	fmt.Println("时间戳转换后的 UTC 时间", t.UTC())
}

func main() {
	TS2Time()
}
