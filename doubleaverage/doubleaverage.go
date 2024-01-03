package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	for i := 5; i > 0; i-- {
		amount := DoubleAverage(5, 100)
		fmt.Println(amount)
	}
}

// 利用二倍均值算法发红包
func DoubleAverage(count, amount int64) int64 {
	const min = 1
	if count == 1 {
		return amount
	}
	max := amount - min*count
	avg := max / count
	avg2 := 2*avg + min
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(avg2) + min

	return x
}
