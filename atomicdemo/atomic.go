package main

import (
	"fmt"
	"math"
	"sync/atomic"
)

func main() {
	// adduint32 在溢出后会自动变成0,不会抛出异常
	var v uint32 = math.MaxUint32 - 1
	atomic.AddUint32(&v, 1)
	fmt.Println(v)
}
