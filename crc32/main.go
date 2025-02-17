package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	cuid := "8CDFE4E11BD78190CBC867D0B25FE4F3|0"
	//cuid := "B5F386F7D4ED912331263B913429A196|0"

	var TableMaxNum int64 = 100
	num := crc32.ChecksumIEEE([]byte(cuid))
	fmt.Println(num)
	suffix := int64(num) % TableMaxNum
	table := fmt.Sprintf("%s%d", "tblUserBuyInfo", suffix)
	fmt.Println(table)
}
