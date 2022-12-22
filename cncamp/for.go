package main

import "fmt"

func main() {
	me := []string{"I", "am", "stupid", "and", "weak"}

	for idx, _ := range me {
		if idx == 2 {
			me[idx] = "smart"
		}
		if idx == 4 {
			me[idx] = "strong"
		}
	}

	fmt.Println(me)
}
