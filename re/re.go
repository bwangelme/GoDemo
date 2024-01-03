package main

import (
	"fmt"
	"regexp"
)

func ExpandStringDemo() {
	content := `
	# comment line
	option1: value1
	option2: value2

	# another comment line
	option3: value3
`

	// https://docs.python.org/zh-cn/3/library/re.html#regular-expression-syntax
	// ?m 匹配一个空字符串，m 表示是多行模式
	// ?m 设置以后，样式字符 '^' 匹配字符串的开始，和每一行的开始（换行符后面紧跟的符号）；
	// 样式字符 '$' 匹配字符串尾，和每一行的结尾（换行符前面那个符号）。
	pattern := regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)$`)

	template := "$key=$value\n"

	result := []byte{}

	for _, submatches := range pattern.FindAllStringSubmatchIndex(content, -1) {
		fmt.Println(submatches)
		// ExpandString 根据 pattern，从 content 中捕获内容
		// 然后用捕获的内容渲染 template
		// 将渲染结果存储到 result 中
		result = pattern.ExpandString(result, template, content, submatches)
	}
	fmt.Println(string(result))
}

func FindMatchDemo() {
	pat := regexp.MustCompile(`^(\d{4}(-\d{1,2}){0,2}).?(\((.*)\))?$`)

	str := "2003-10-19日"
	paras := pat.FindStringSubmatch(str)
	fmt.Println(paras[1])
}

func main() {
	FindMatchDemo()
	//ExpandStringDemo()
}