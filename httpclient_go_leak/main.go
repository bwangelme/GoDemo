package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func Index(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 100; i++ {
		go func() {
			http.Get("https://www.baidu.com/")
		}()
	}
	fmt.Fprintf(w, "ok")
}

func main() {
	http.HandleFunc("/", Index)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
