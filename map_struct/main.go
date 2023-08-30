/*
range 循环 map 的时候 ，只写一个接收者，它就是 key
*/
package main

import "fmt"

func main() {
	v := map[string]struct{}{}

	v["key1"] = struct{}{}
	v["key2"] = struct{}{}
	v["key3"] = struct{}{}

	for key := range v {
		fmt.Println(key)
	}
}
