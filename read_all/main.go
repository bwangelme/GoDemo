package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	var buf = make([]byte, 4)
	r1 := strings.NewReader("abcd")
	r2 := strings.NewReader("1234")

	n, err := io.ReadFull(r1, buf)
	fmt.Println(n, err)
	fmt.Println(string(buf))
	n, err = io.ReadFull(r2, buf)
	fmt.Println(n, err)

	fmt.Println(string(buf))
}
