package main

import "fmt"

type Span struct {
	id int64
}


func main() {
	var a []*Span
	fmt.Println(a)
	fmt.Println(append(a, &Span{2}))
	fmt.Println(append(a, nil))
	fmt.Println(a)
}
