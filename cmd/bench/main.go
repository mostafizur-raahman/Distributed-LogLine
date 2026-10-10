package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var (
		url         string
		total       int
		concurrency int
	)

	flag.StringVar(&url, "url", "http://localhost:4000/ingest", "target URL")
	flag.IntVar(&total, "n", 5000, "total requests")
	flag.IntVar(&concurrency, "c", 50, "concurrent workers")
	flag.Parse()

	if concurrency > total {
		concurrency = total
	}

	fmt.Printf("bench: %d requests, %d concurrent → %s\n\n", total, concurrency, url)

	payload := []byte(`{"level":"info","message":"bench entry","service":"bench","timestamp":"2026-03-01T00:00:00Z"}`)

	var (
		success   atomic.Int64
		failed    atomic.Int64
		mu        sync.Mutex
		latencies []time.Duration
	)

	latencies = make([]time.Duration, 0, total)
	work := make(chan int, total)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        concurrency,
			MaxIdleConnsPerHost: concurrency,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	var wg sync.WaitGroup

	start := time.Now()

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range work {
				t := time.Now()
				resp, err := client.Post(url, "application/json", bytes.NewReader(payload))
				d := time.Since(t)

				if err != nil {
					failed.Add(1)
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == http.StatusCreated {
					success.Add(1)
				} else {
					failed.Add(1)
				}

				mu.Lock()
				latencies = append(latencies, d)
				mu.Unlock()
			}
		}()
	}

	for i := 0; i < total; i++ {
		work <- i
	}
	close(work)
	wg.Wait()

	elapsed := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	fmt.Printf("completed in %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  total:    %d\n", total)
	fmt.Printf("  success:  %d\n", success.Load())
	fmt.Printf("  failed:   %d\n", failed.Load())
	fmt.Printf("  rps:      %.0f\n", float64(success.Load())/elapsed.Seconds())

	if len(latencies) > 0 {
		fmt.Printf("  p50:      %s\n", percentile(latencies, 0.50))
		fmt.Printf("  p95:      %s\n", percentile(latencies, 0.95))
		fmt.Printf("  p99:      %s\n", percentile(latencies, 0.99))
		fmt.Printf("  max:      %s\n", latencies[len(latencies)-1])
	}
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	idx := int(float64(len(sorted)) * p)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
