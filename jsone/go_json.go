package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type QuoteString string

func (q QuoteString) MarshalJSON() ([]byte, error) {
	data := []byte(strconv.QuoteToASCII(string(q)))
	return data, nil
}

type Message struct {
	ID   int         `json:"id"`
	Name QuoteString `json:"name"`
}

// 使用自定义的 Marshal 方法来数据 JSON 数据
func testCustomMarshal() {
	var m = Message{
		Name: "golang中文Unicode编码",
		ID:   1,
	}
	res, err := json.Marshal(m)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(res))
}

type S3 struct {
	A json.Number `json:"a"`
}

// 使用 json.Number 类型来解析数字
func testNumber() {
	var raw = []byte(`{"a":1}`)
	var raw_string = []byte(`{"a":"2"}`)

	var s S3
	err := json.Unmarshal(raw, &s)
	fmt.Println(err, s.A, reflect.TypeOf(s.A))

	err = json.Unmarshal(raw_string, &s)
	fmt.Println(err, s.A, reflect.TypeOf(s.A))
}

// 使用 useNumber 选项来自动解析数字
func test3() {
	decoder := json.NewDecoder(bytes.NewBufferString(`{"10000000000":10000000000,"111":1}`))
	decoder.UseNumber()
	var obj map[string]interface{}
	decoder.Decode(&obj)
	v := obj["10000000000"]
	rV := v.(json.Number)
	rrV, _ := rV.Int64()
	fmt.Println(reflect.TypeOf(v), rrV)
}

func main() {
	testCustomMarshal()
	testNumber()
	test3()
}
