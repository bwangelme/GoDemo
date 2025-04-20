package main

func main() {
	var m = make(map[int]string, 0)

	// 给 map 中的数据赋值
	m[1] = "hello"

	// 访问 map 中的数据
	_ = m[2]

	// 遍历 map
	for range m {
	}
}
