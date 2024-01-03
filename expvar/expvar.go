package main

import (
	"expvar"
	"fmt"
	"net/http"
)

var visits = expvar.NewInt("visits")

func handler(w http.ResponseWriter, r *http.Request) {
	visits.Add(1)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	// expvar 包提供了一些接口，可以监听一些变量，然后通过 /debug/var http 接口访问到
	// 看起来就很像是监控用的，prometheus 的监控模式就是这种，定期从接口中获取指标
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
