package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		PoolSize: 40,
	})

	rdb.Set(ctx, "First value", "value_1", 0)
	rdb.Set(ctx, "Second value", "value_2", 0)

	var group sync.WaitGroup

	for i := 0; i < 20; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			val := rdb.Get(ctx, "Second value").Val()
			if val != "value_2" {
				log.Fatalf("val was not set. expected: %s but got: %s", "value_2", val)
			} else {
				log.Printf("Test Success")
			}
		}()
	}
	group.Wait()

	rdb.Del(ctx, "First value")
	rdb.Del(ctx, "Second value")
}
