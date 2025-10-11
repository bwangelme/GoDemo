package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/trace"
	"strconv"
	"sync"
	"time"
)

// bucket is a simple mutex-protected counter.
type bucket struct {
	mu      sync.Mutex
	guesses int
}

func main() {
	// Set up the flight recorder
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   1000 * time.Millisecond,
		MaxBytes: 1 << 22, // 1 MiB
	})
	fr.Start()

	// Make one bucket for each valid number a client could guess.
	// The HTTP handler will look up the guessed number in buckets by
	// using the number as an index into the slice.
	buckets := make([]bucket, 100)

	// Every minute, we send a report of how many times each number was guessed.
	go func() {
		for range time.Tick(2 * time.Second) {
			sendReport(buckets)
		}
	}()

	// Choose the number to be guessed.
	answer := rand.Intn(len(buckets))

	http.HandleFunc("/guess-number", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Fetch the number from the URL query variable "guess" and convert it
		// to an integer. Then, validate it.
		guess, err := strconv.Atoi(r.URL.Query().Get("guess"))
		if err != nil || !(0 <= guess && guess < len(buckets)) {
			http.Error(w, "invalid 'guess' value", http.StatusBadRequest)
			return
		}

		// Select the appropriate bucket and safely increment its value.
		b := &buckets[guess]
		b.mu.Lock()
		b.guesses++
		b.mu.Unlock()

		// Respond to the client with the guess and whether it was correct.
		fmt.Fprintf(w, "guess: %d, correct: %t", guess, guess == answer)

		// Capture a snapshot if the response takes more than 100ms.
		// Only the first call has any effect.
		if fr.Enabled() && time.Since(start) > 500*time.Millisecond {
			go captureSnapshot(fr)
		}

		log.Printf("HTTP request: endpoint=/guess-number guess=%d duration=%s", guess, time.Since(start))
	})
	log.Fatal(http.ListenAndServe(":8090", nil))
}

// sendReport posts the current state of buckets to a remote service.
func sendReport(buckets []bucket) {
	counts := make([]int, len(buckets))

	for index := range buckets {
		b := &buckets[index]
		b.mu.Lock()
		defer b.mu.Unlock()

		counts[index] = b.guesses
	}

	// Marshal the report data into a JSON payload.
	b, err := json.Marshal(counts)
	if err != nil {
		log.Printf("failed to marshal report data: error=%s", err)
		return
	}
	url := "http://localhost:8091/guess-number-report"
	if _, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
		log.Printf("failed to send report: %s", err)
	}
}

var once sync.Once

// captureSnapshot captures a flight recorder snapshot.
func captureSnapshot(fr *trace.FlightRecorder) {
	// once.Do ensures that the provided function is executed only once.
	once.Do(func() {
		f, err := os.Create("snapshot.trace")
		if err != nil {
			log.Printf("opening snapshot file %s failed: %s", f.Name(), err)
			return
		}
		defer f.Close() // ignore error

		// WriteTo writes the flight recorder data to the provided io.Writer.
		_, err = fr.WriteTo(f)
		if err != nil {
			log.Printf("writing snapshot to file %s failed: %s", f.Name(), err)
			return
		}

		// Stop the flight recorder after the snapshot has been taken.
		fr.Stop()
		log.Printf("captured a flight recorder snapshot to %s", f.Name())
	})
}
