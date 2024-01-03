package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.MkdirAll("./data", os.ModePerm)
	fmt.Println(err)
}
