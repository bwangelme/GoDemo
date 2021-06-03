package main

// source: https://v2ex.com/t/780962#reply24

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func checkI(i int) bool {
	fmt.Println("CheckI", i)
	return i < 10
}

func testFor() {
	var i int
	for ; checkI(i); i++ {
		// break 之后不会再去执行 ++ 操作和条件检查
		if i == 0 {
			i++
			break
		}
	}
	fmt.Println(i)
}

type Token struct {
	Data  map[string]string `json:"data"`
	Order []string          `json:"order"`
}

func (t *Token) Value() (res string) {
	for _, key := range t.Order {
		val := t.Data[key]
		res += val
	}

	return res
}

func (t *Token) validate() error {
	for _, key := range t.Order {
		_, ok := t.Data[key]
		if !ok {
			return fmt.Errorf("key %v not exist in order", key)
		}
	}

	return nil
}

func NewToken(jsonbyte []byte) (*Token, error) {
	t := &Token{}
	err := json.Unmarshal(jsonbyte, &t)
	if err != nil {
		return nil, err
	}

	err = t.validate()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func testJSONUnmarshal() {
	jsonStr := `{"data":{"name":"tom","user_id":"123"}, "order": ["user_id", "name"]}`
	token, err := NewToken([]byte(jsonStr))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(token.Value())
}

func testGJson() {
	jsonStr := `{"name":"tom","user_id":"123"}`
	r := gjson.Parse(jsonStr)
	fmt.Println(r.Raw)
	// gjson 执行 foreach 的时候，是按照 json 字符串中的顺序进行 unmarshal 的
	r.ForEach(func(key, value gjson.Result) bool {
		fmt.Println(key, value)
		return true
	})
}

func main() {
	testJSONUnmarshal()
}
