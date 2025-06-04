package main

import (
	"context"
	"fmt"
)

/**
字节面试题

本程序测试如果存在多个同名 key 的 value, 最新设置的会被获取到

因为多个 ctx 组成一个链表，查找时从链表尾部开始查找，所以最新设置的会被获取到
*/

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key", "v1")
	ctx = context.WithValue(ctx, "key", "v2")

	fmt.Println(ctx.Value("key"))
}
