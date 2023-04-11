package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Resp struct {
	Code int      `json:"code,omitempty"`
	Urls []string `json:"urls,omitempty"`
}

func main() {
	var resp = Resp{}
	//var body = []byte(`{"name": "bwangel"}`)
	//err := json.Unmarshal(body, &resp)

	var r = strings.NewReader(`{"name": "bwangel"}`)
	deco := json.NewDecoder(r)

	deco.DisallowUnknownFields()
	err := deco.Decode(&resp)
	fmt.Println(resp, err)

}
