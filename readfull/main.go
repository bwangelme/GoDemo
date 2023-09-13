package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("abc")
	var buf = make([]byte, 4)

	for {
		n, err := io.ReadFull(r, buf)
		fmt.Println(n, err, len(buf))
		break
	}
}
