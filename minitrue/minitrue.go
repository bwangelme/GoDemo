package main

import (
	"fmt"
	"net/http"
	"reflect"
)

// Cond
// 添加运算符的简单实现，支持任意类型
func Cond[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}

// Or
// 返回 vals 中第一个不是对应类型0值的参数，如果没有，则返回对应类型的0值
func Or[T comparable](vals ...T) T {
	for _, val := range vals {
		if val != *new(T) {
			return val
		}
	}
	return *new(T)
}

func main() {
	var req1 = &http.Request{
		Method: http.MethodPost,
	}
	var req *http.Request
	res := Cond(req == nil, http.MethodGet, http.MethodPost)
	fmt.Println(res, reflect.TypeOf(res))

	a := Or(req, req1)
	fmt.Println(a.Method)
}
