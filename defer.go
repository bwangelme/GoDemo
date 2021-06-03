package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("Defer")
	}()

	fmt.Println("YES")
}
