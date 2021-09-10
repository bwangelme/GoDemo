package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type ContextKey string

const ContextKeyTrace ContextKey = "trace"

type Trace struct {
	id    int64
	spans []string
}

func NewTrace() *Trace {
	return &Trace{
		id: rand.Int63n(2048),
	}
}

func (t *Trace) Submit() {
	t.spans = append(t.spans, fmt.Sprintf("span %d", rand.Intn(1024)))
	fmt.Println(t.spans)
}

func (t Trace) String() string {
	return fmt.Sprintf("Trace<%v>", t.id)
}

func NewContext(ctx context.Context, trace *Trace) context.Context {
	ctx = context.WithValue(ctx, ContextKeyTrace, trace)
	return ctx
}

func FromContext(ctx context.Context) (*Trace, error) {
	rawValue := ctx.Value(ContextKeyTrace)

	trace, ok := rawValue.(*Trace)
	if !ok {
		return nil, fmt.Errorf("invalid trace %v", rawValue)
	}

	return trace, nil
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := NewTrace()
		r = r.WithContext(NewContext(ctx, t))

		next.ServeHTTP(w, r)
	})
}

func worker(i int, r *http.Request) {
	time.Sleep(2 * time.Second)
	runtime.GC()
	fmt.Println("Start Goroutine", i)
	trace, err := FromContext(r.Context())
	if err != nil {
		log.Panicln(err)
	}
	trace.Submit()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 3; i++ {
		go worker(i, r)
	}

	fmt.Fprintf(w, "hello, world")
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(HomeHandler))

	app := TraceMiddleware(mux)

	http.ListenAndServe(":8080", app)
}
