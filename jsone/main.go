package main

//将 json 内容转换成字符串，再执行 json encode, 输出的文本中，引号前就会多一个 \ 符号
// {"content":"{\"code\":1,\"data\":{\"adId\":\"5005629286175060\",\"adLogo\":\"\"}}"}

import (
	"encoding/json"
	"fmt"
)

func main() {
	res := make(map[string]interface{})
	res["code"] = 1
	res["data"] = map[string]string{
		"adId":   "5005629286175060",
		"adLogo": "",
	}

	tLog, _ := json.Marshal(res)
	ttLog, _ := json.Marshal(map[string]string{
		"content": string(tLog),
	})
	fmt.Printf("%s", fmt.Sprintf("%s", string(ttLog)))
}
