package main

import "fmt"

type RouterGroup struct {
}

func (group *RouterGroup) group() {
	fmt.Println("group")
}

type Engine struct {
	RouterGroup
}

func main() {
	// Engine 对象能够直接调用 RouterGroup 对象的方法
	e := Engine{RouterGroup{}}
	e.group()
}
