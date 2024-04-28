package main

import (
	"bwdemo/reqcache/middleware"
	"fmt"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
)

func Home(w http.ResponseWriter, r *http.Request) {
	reqid, _ := middleware.FromContext(r.Context())
	w.WriteHeader(http.StatusOK)

	c, _ := middleware.CacheFromContext(r.Context())
	c.Add("some_key", "some val")
	v, ok := c.Get("some_key")

	fmt.Fprintf(w, "Req id: %v, v: %v|%v\n", reqid, v, ok)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.Handle("/_dae/debug/pprof/", http.StripPrefix("/_dae", http.HandlerFunc(pprof.Index)))
	mux.Handle("/_dae/debug/pprof/cmdline", http.StripPrefix("/_dae", http.HandlerFunc(pprof.Cmdline)))
	mux.Handle("/_dae/debug/pprof/profile", http.StripPrefix("/_dae", http.HandlerFunc(pprof.Profile)))
	mux.Handle("/_dae/debug/pprof/symbol", http.StripPrefix("/_dae", http.HandlerFunc(pprof.Symbol)))
	mux.Handle("/_dae/debug/pprof/trace", http.StripPrefix("/_dae", http.HandlerFunc(pprof.Trace)))

	m := middleware.ReqIDMiddleware(mux)
	m = middleware.ReqCacheMiddleware(m)
	// Start the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", m)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
