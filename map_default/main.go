package main

import "fmt"

// ParseRefererChannel
// 测试 map 在获取不存在的值的时候，不会 panic
func ParseRefererChannel(refererCh string) int {
	// default is 0
	var refererChannel = 0
	refererMap := map[string]int{
		"picSearchPvalnlabel": 1, // 快搜结果页的非准召结果的提问入口
	}
	refererChannel, _ = refererMap[refererCh]

	return refererChannel
}

func main() {
	fmt.Println(ParseRefererChannel("default"))             // 0
	fmt.Println(ParseRefererChannel("picSearchPvalnlabel")) // 1
	fmt.Println(ParseRefererChannel(""))                    // 0
}
