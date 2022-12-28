package main

import "fmt"

type a int

func (a) A() {
	fmt.Println("a")
}

type http struct {
}

func (h *http) Get() {
	fmt.Println("Get")
}

func main() {
	var methd = interface{ Get() }.Get
	var hh http
	methd(&hh)
}
