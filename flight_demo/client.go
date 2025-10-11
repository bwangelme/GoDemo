package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Define command-line flags to configure the load test.
	concurrency := flag.Int("concurrency", 100, "number of concurrent workers")
	calls := flag.Int("calls", 10000, "total number of calls")
	serverAddr := flag.String("server-addr", "localhost:8090", "server address")
	flag.Parse()

	// Validate the provided flags.
	if *concurrency < 1 || *calls < 1 || *serverAddr == "" {
		log.Fatalf("invalid argument value")
	}

	log.Printf("client: starting %d concurrency, %d total-calls", *concurrency, *calls)
	var wg sync.WaitGroup

	// Calculate the number of calls for each worker.
	perGCallCount := *calls / *concurrency
	// The remainder is added to the last worker to ensure the total number of calls is met.
	callCountRemainder := *calls % *concurrency

	// Start the worker goroutines.
	for workerID := 0; workerID < *concurrency; workerID++ {
		// wg.Go starts a goroutine and registers the call with the WaitGroup.
		wg.Go(func() {
			// Each goroutine gets its own HTTP client.
			c := http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 5,
				},
			}
			callCount := perGCallCount
			// The first worker also handles the remainder of the calls.
			if workerID == 0 {
				callCount = callCount + callCountRemainder
			}

			// Each worker makes its assigned number of HTTP calls.
			for _ = range callCount {
				// Generate a random number for the guess.
				url := fmt.Sprintf("http://%s/guess-number?guess=%d", *serverAddr, rand.Intn(100))

				// Make the HTTP GET request.
				resp, err := c.Get(url)
				if err != nil {
					log.Printf("client: failed to retrieve %s: %s", url, err)
					continue
				}
				// Discard the response body to allow the connection to be reused.
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		})
	}

	s := &http.Server{
		Addr:           ":8091",
		Handler:        &reportHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	// Wait for all goroutines to complete.
	wg.Wait()
	log.Printf("client: DONE %d concurrency, %d total-calls", *concurrency, *calls)
	s.Shutdown(context.Background())
	log.Printf("client: HTTP server shutdown complete")
}

type reportHandler struct{}

func (rh *reportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sleepDuration := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepDuration)
	fmt.Fprint(w, "Report recieved")
}
