package main

import (
	"fmt"
	"regexp"
)

func main() {
	pat := regexp.MustCompile(`^(\d{4}(-\d{1,2}){0,2}).?(\((.*)\))?$`)

	str := "2003-10-19日"
	paras := pat.FindStringSubmatch(str)
	fmt.Println(paras[1])
}
