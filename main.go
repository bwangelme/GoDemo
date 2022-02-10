package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	inner := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:   0, // use default DB
		})
	i, _ := inner.Info(context.Background()).Result()
	fmt.Println(i)
}
