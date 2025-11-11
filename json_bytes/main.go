package main

import (
	"encoding/json"
	"fmt"
)

// golang 中，对 bytes 进行 json 序列化的话，golang 会对它进行 base64 编码，输出一个 " 包裹的 base64 string
func main() {
	fmt.Println("=== 嵌套序列化：第一次序列化 -> bytes -> 第二次序列化 ===")

	// 第一层：内部数据结构
	type InnerData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	inner := InnerData{
		Name: "张三",
		Age:  25,
	}

	// 第一次序列化：将 InnerData 序列化为 JSON
	jsonBytes, _ := json.Marshal(inner)
	fmt.Printf("第一次序列化结果 ([]byte):\n%v\n\n", jsonBytes)
	fmt.Printf("第一次序列化结果 (string):\n%v\n\n", string(jsonBytes))

	// 第二次序列化：直接对 []byte 进行序列化
	// Go 的 json.Marshal 对 []byte 类型会进行 base64 编码，然后作为 JSON 字符串输出
	// 因此结果是一个用双引号包裹的字符串
	finalJSON, _ := json.Marshal(jsonBytes)
	fmt.Printf("第二次序列化结果([]byte):\n%v \"=%d e=%d\n\n", finalJSON, finalJSON[0], finalJSON[1])
	fmt.Printf("第二次序列化结果(string):\n%s\n\n", string(finalJSON))
}
