package main

import (
	"bwdemo/pkg"
	"fmt"
	"time"
)

// 五个协程打印的 client.ID 是相同的，说明文件中的全局变量只声明一次
func main() {
	for i := 0; i < 5; i++ {
		go func(i int) {
			c := pkg.DefaultClient
			fmt.Printf("%d goroutine, client id: %v\n", i, c.ID)
		}(i)
	}

	time.Sleep(time.Second * 2)
}
