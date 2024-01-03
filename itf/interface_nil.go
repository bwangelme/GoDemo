package main

import "fmt"

type User struct {
	Name string
}

func Has(user interface{}) {
	fmt.Println(user == nil)
}

func main() {
	var u *User = nil
	fmt.Println(u == nil)
	// 传入的是 *User 类型的 nil，所以输出 false
	Has(u)   // false
	Has(nil) // true
}
