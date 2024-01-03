package main

import "fmt"

func main() {
	var data = make(map[string]interface{})
	var metrics, ok = data["metrics"].([]map[string]interface{})

	fmt.Println(metrics, ok)

}
