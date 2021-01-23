package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Config struct {
	a []int
}

func (c *Config) T() {}

// 产生数据竞争的代码，它的行为是很难预测的
// 不要研究存在数据竞争的代码的行为
func main() {
	// 协程们会在 cfg 这个变量上产生数据竞争
	cfg := &Config{}

	go func() {
		i := 0
		for {
			i++
			cfg.a = []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < runtime.NumCPU(); n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 100; n++ {
				fmt.Println(cfg.a)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
