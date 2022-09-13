package main

import (
	"bwdemo/protojson/apppb"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	var app = apppb.App{}
	jsonData := []byte(`{"id": "10"}`)
	err := protojson.Unmarshal(jsonData, &app)
	fmt.Println(err)

	fmt.Println(app.Id)
}
