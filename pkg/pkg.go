package pkg

import "fmt"

var (
	Key = "value"
)

func init() {
	fmt.Println("YES")
	Key = "val"
}

func SomeAction() {
	fmt.Println("Action")
}
