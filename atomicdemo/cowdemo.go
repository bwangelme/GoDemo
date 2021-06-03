package main

import (
	"sync"
	"sync/atomic"
)

// Go 业务中的 Copy On Write
// 当有配置需要更改时，先复制一份，在复制份上执行修改，然后再利用 atomic.Value 这个原子操作修改配置的指针
func main() {
	type Map map[string]string
	var m atomic.Value
	var mu sync.Mutex // 写协程之间的互斥锁

	read := func(key string) (val string) {
		m1 := m.Load().(Map)
		return m1[key]
	}
	write := func(key, value string) {
		mu.Lock()
		defer mu.Unlock()

		m1 := m.Load().(Map)
		m2 := make(Map)
		for k, v := range m1 {
			m2[k] = v
		}
		m2[key] = value
		m.Store(m2)
	}

	_, _ = read, write
}
