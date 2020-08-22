package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type TestRedisFunc func(rdb *redis.Client, ctx context.Context)

var benchTimes = 10000

func PingWithPipelining(rdb *redis.Client, ctx context.Context) {
	_, err := rdb.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for i := 0; i < benchTimes; i++ {
			pipeliner.Ping(ctx)
		}
		return nil
	})
	if err != nil {
		log.Printf("%v\n", err)
	}
}

func PingWithoutPipelining(rdb *redis.Client, ctx context.Context) {
	for i := 0; i < benchTimes; i++ {
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Printf("%d: %v\n", i, err)
		}
	}
}

// 执行性能测试
func Benchmark(callback TestRedisFunc) time.Duration {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	start := time.Now()
	callback(rdb, ctx)
	return time.Now().Sub(start)
}

func main() {

	for _, tt := range []struct {
		Name     string
		Callback TestRedisFunc
	}{
		{
			Name:     "WithPipelining",
			Callback: PingWithPipelining,
		},
		{
			Name:     "WithoutPipelining",
			Callback: PingWithoutPipelining,
		},
	} {
		fmt.Printf("Start %v\n", tt.Name)
		duration := Benchmark(tt.Callback)
		fmt.Printf("%v Running %d times cost %v\n", tt.Name, benchTimes, duration)
	}
}
