package main

import "fmt"

func f(v *[][]int) {
	*v = append(*v, make([]int, 3))
	*v = append(*v, make([]int, 3))

	(*v)[1][2] = 42
}

func main() {
	var res = make([][]int, 0)

	f(&res)

	fmt.Println(res)
}
