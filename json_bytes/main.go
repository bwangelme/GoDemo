package main

import (
	"encoding/json"
	"fmt"
)

// golang 中，对 bytes 进行 json 序列化的话，golang 会对它进行 base64 编码
// 本程序的输出如下
// 第一次序列化结果 ([]byte):
// [123 34 110 97 109 101 34 58 34 229 188 160 228 184 137 34 44 34 97 103 101 34 58 50 53 125]
//
// 第二次序列化结果:
// "eyJuYW1lIjoi5byg5LiJIiwiYWdlIjoyNX0="
func main() {
	fmt.Println("=== 嵌套序列化：第一次序列化 -> 转string -> 第二次序列化 ===\n")

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

	// 第二次序列化：直接对 []byte 进行序列化
	finalJSON, _ := json.Marshal(jsonBytes)
	fmt.Printf("第二次序列化结果:\n%s\n\n", string(finalJSON))
}
