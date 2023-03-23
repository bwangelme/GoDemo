package main

import "fmt"

type student struct {
	name string
	age int
}

func (s student) String() string {
	return fmt.Sprintf("%v: %v", s.name, s.age)
}

func main() {
	m := make(map[string]*student)
	stuts := []student{
		{name: "小王子", age: 21},
		{name: "娜扎", age: 22},
		{name: "大王八", age: 23},
	}

	// for range 遍历时用的是同一个地址空间来存储当前遍历到的 student 结构，
	// 因为使用的是指针，后面的循环把前面的值给改了
	// 所以运行结束后打印出来都是最后一次循环遍历到的 student 内容。
	for _, stu := range stuts {
		// 正确做法是将 stu 存到一个新的空间中
		// s := stu
		m[stu.name] = &stu
		fmt.Println(stu)
	}

	for k, v := range m {
		fmt.Printf("%v => %v: %v\n", k, v.name, v.age)
	}
}
