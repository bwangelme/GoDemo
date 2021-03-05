package main

import (
	"fmt"
	"log"
	"unicode"
)

const (
	MAX_KEY_LEN = 250
)

func IsValidKeyString(key string) bool {
	length := len(key)
	if length == 0 || length > MAX_KEY_LEN {
		log.Fatal("bad key len=%d", length)
		return false
	}

	if key[0] <= ' ' || key[0] == '?' || key[0] == '@' {
		log.Fatal("bad key len=%d key[0]=%x", length, key[0])
		return false
	}

	for _, r := range key {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			log.Fatal("bad key len=%d %s", length, key)
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(IsValidKeyString("中文"))
	fmt.Println(IsValidKeyString("汉字"))
}
