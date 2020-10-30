package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	secret := "SGVsbG8sIEpldGJyYWlucyBpbiAxMDI0Cg=="
	data, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s", data)
}
