package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
)

func main() {
	err := hystrix.Do("get_baidu", func() error {
		// talk to other services
		_, err := http.Get("https://www.baidu.com/")
		if err != nil {
			fmt.Println("get error")
			return err
		}
		return nil
	}, func(err error) error {
		fmt.Println("get an error, handle it")
		return nil
	})

	if err != nil {
		fmt.Println("get_baidu error", err)
	}
}
