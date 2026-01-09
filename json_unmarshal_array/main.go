package main

import (
	"encoding/json"
	"fmt"
)

// 定义一个包含数组的 struct
type Data struct {
	Name   string   `json:"name"`
	Items  []string `json:"items"`
	Numbers []int   `json:"numbers"`
}

func main() {
	fmt.Println("=== 测试 json.Unmarshal 两次解析不同数组的行为 ===\n")

	// 创建同一个变量
	var data Data

	// 第一次解析：包含数组的 JSON
	jsonStr1 := `{
		"name": "第一次",
		"items": ["apple", "banana", "cherry"],
		"numbers": [1, 2, 3]
	}`

	fmt.Println("第一次解析的 JSON:")
	fmt.Println(jsonStr1)
	fmt.Println()

	err := json.Unmarshal([]byte(jsonStr1), &data)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	fmt.Println("第一次解析后的结果:")
	fmt.Printf("  Name: %s\n", data.Name)
	fmt.Printf("  Items: %v\n", data.Items)
	fmt.Printf("  Numbers: %v\n", data.Numbers)
	fmt.Printf("  Items 长度: %d\n", len(data.Items))
	fmt.Printf("  Numbers 长度: %d\n", len(data.Numbers))
	fmt.Println()

	// 第二次解析：不同的数组内容
	jsonStr2 := `{
		"name": "第二次",
		"items": ["dog", "cat"],
		"numbers": [10, 20, 30, 40, 50]
	}`

	fmt.Println("第二次解析的 JSON:")
	fmt.Println(jsonStr2)
	fmt.Println()

	err = json.Unmarshal([]byte(jsonStr2), &data)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	fmt.Println("第二次解析后的结果:")
	fmt.Printf("  Name: %s\n", data.Name)
	fmt.Printf("  Items: %v\n", data.Items)
	fmt.Printf("  Numbers: %v\n", data.Numbers)
	fmt.Printf("  Items 长度: %d\n", len(data.Items))
	fmt.Printf("  Numbers 长度: %d\n", len(data.Numbers))
	fmt.Println()

	// 结论
	fmt.Println("=== 结论 ===")
	fmt.Println("✓ 数组会被完全刷新（替换），而不是追加")
	fmt.Println("✓ 第一次解析的数组内容不会保留")
	fmt.Println("✓ 第二次解析会完全覆盖第一次的结果")
	fmt.Println("✓ 如果第二次解析的数组更短，数组长度会变短")
	fmt.Println("✓ 如果第二次解析的数组更长，数组长度会变长")
}

