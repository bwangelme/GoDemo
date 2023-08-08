package main

/*
本程序测试 http 加上中间件后，client 收到请求的返回时间

中间件中调用了 next.ServeHTTP 之后，client 并不会立刻收到请求，只有所有中间件都执行结束之后，client 才会收到请求
*/

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewLongtimMiddleware(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(writer, request)
		time.Sleep(10 * time.Second)
	}
	return http.HandlerFunc(fn)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "ok\n")
	})

	err := http.ListenAndServe("localhost:8080", NewLongtimMiddleware(mux))
	if err != nil {
		log.Fatalln(err)
	}
}
