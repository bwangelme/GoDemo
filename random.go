package main

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
)

func main() {
	// 伪随机数
	rand.Seed(2)
	for i := 0; i < 4; i++ {
		fmt.Println(rand.Intn(100))
	}

	fmt.Println(strings.Repeat("-", 16))
	// 真正的随机
	for i := 0; i < 4; i++ {
		n, _ := crand.Int(crand.Reader, big.NewInt(100))
		fmt.Println(n.Int64())
	}
}
