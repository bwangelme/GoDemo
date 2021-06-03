package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// atomic 比 rwlock 要快一个数量级
//ø> go test -bench=. atomicdemo/config_test.go                                                                                                                                                                      21:49:13 (01-21)
//goos: linux
//goarch: amd64
//BenchmarkMutex-16       1000000000               0.000280 ns/op
//BenchmarkAtomic-16      1000000000               0.000061 ns/op
//PASS
//ok      command-line-arguments  0.013s

type Config struct {
	a []int
}

func (c *Config) T() {}

func BenchmarkMutex(b *testing.B) {
	var l sync.RWMutex
	var cfg *Config

	go func() {
		i := 0
		for {
			i++
			l.Lock()
			cfg = &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			l.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 100; n++ {
				l.RLock()
				cfg.T()
				l.RUnlock()
				//fmt.Println(cfg)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkAtomic(b *testing.B) {
	var v atomic.Value
	v.Store(&Config{})

	go func() {
		i := 0
		for {
			i++
			cfg := &Config{a: []int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			v.Store(cfg)
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < 100; n++ {
				cfg := v.Load().(*Config)
				cfg.T()
				//fmt.Println(cfg)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
