package main

func main() {
	var m = map[int]string{
		3: "world",
		4: "!",
	}

	// 给 map 中的数据赋值
	m[1] = "hello"

	// 访问 map 中的数据
	_ = m[2]

	// 遍历 map
	for range m {
	}

	delete(m, 4)
}
