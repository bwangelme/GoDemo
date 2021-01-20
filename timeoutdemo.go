package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func search(term string) (string, error) {
	time.Sleep(time.Millisecond * 200)
	return "some value", nil
}

type result struct {
	record string
	err    error
}

func Process() error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := make(chan result)

	// 并发操作由调用者来执行
	go func() {
		record, err := search("中文")
		ch <- result{record: record, err: err}
	}()

	select {
	case result := <-ch:
		if result.err != nil {
			return result.err
		}
		fmt.Println("Received:", result.record)
		return nil
	case <-ctx.Done():
		return errors.New("search canceled")
	}
}

func main() {
	err := Process()
	if err != nil {
		log.Printf("Process err: %v", err)
	}
}
