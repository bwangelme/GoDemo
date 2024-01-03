package main

import (
	"fmt"
	"sort"
)

func mapInit() {
	// map 中不存在的key会自动初始化成0值，可以执行 += 操作
	currentJobs := make(map[string]int)
	for i := 0; i < 10; i++ {
		currentJobs[string(rune('a'+i))] += i
	}
	currentJobs["B"] += 20
	//fmt.Println(v, ok)

	fmt.Println(currentJobs)
}

func mapIter() {
	var weightMap map[string]int = map[string]int{
		"belba1:8000#0": 1,
		"belba3:8000#0": 1,
		"belba2:8000#0": 1,
	}
	addresses := make([]string, 0, len(weightMap))
	for addr := range weightMap {
		fmt.Println(addr)
		addresses = append(addresses, addr)
	}
	fmt.Println(addresses)
	sort.StringSlice(addresses).Sort()
	fmt.Println(addresses)

}

func main() {
	mapIter()
}
