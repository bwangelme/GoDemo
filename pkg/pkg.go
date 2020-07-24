package pkg

import "fmt"

var (
	Key = "value"
)

// 就算没有引用 pkg 文件中的方法，变量，调用其他文件的方法变量时，pkg.go 中的 init 方法也会被执行
func init() {
	fmt.Println("Pkg Init")
	Key = "val"
}

func SomeAction() {
	fmt.Println("Action")
}
