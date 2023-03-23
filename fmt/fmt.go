package main

import (
	"fmt"
)

type OrderBy uint8

const (
	SCORE OrderBy = iota
	RatingCount
	CollectCount
)

// String 方法应该使用值接受器
func (o OrderBy) String() string {
	if o == CollectCount {
		return "collect_count"
	} else if o == RatingCount {
		return "rating_count"
	} else {
		return "score"
	}
}

func main() {
	// %v 能兼容任意格式
	fmt.Printf("%s, %v, %v\n", "abc", 1, 233)

	fmt.Printf("%v\n", SCORE)
	var po = new(OrderBy)
	*po = CollectCount
	fmt.Println(po)
}
