package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	filename := "/Users/michaeltsui/go/pkg/mod/github.com/golang/protobuf@v1.3.3/go.mod"

	fd, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Fatalln(err)
	}
	sum := sha256.Sum256(data)
	src := sum[:]
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, src)
	fmt.Printf("%s\n", dst)
}
