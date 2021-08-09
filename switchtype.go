package main

import "fmt"

func main() {
	var t interface{} = 1

	switch v := t.(type) {
	case int, string:
		fmt.Println("got", v)
	}
}
