package json_empty

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Resp struct {
	Code int      `json:"code,omitempty"`
	Urls []string `json:"urls,omitempty"`
}

// Marshal_int_str
// 在 json tag 中添加了 string 字段
// 序列化出来的 json，它将 Uid 字段变成了 string 类型
func Marshal_int_str() string {
	var Data = struct {
		Uid int64 `json:"uid,string"`
	}{
		Uid: 12345578,
	}

	v, _ := json.Marshal(Data)
	return string(v)
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
