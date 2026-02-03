package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/wingnut128/outlier-go/internal/config"
	"github.com/wingnut128/outlier-go/pkg/api"
)

// TestStress_ConcurrentAPIRequests tests the server under concurrent load
func TestStress_ConcurrentAPIRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	const (
		numGoroutines        = 50
		requestsPerGoroutine = 20
	)

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*requestsPerGoroutine)
	successCount := int64(0)
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				// Generate random values
				numValues := 100 + rand.Intn(900)
				values := make([]float64, numValues)
				for k := 0; k < numValues; k++ {
					values[k] = rand.Float64() * 1000
				}

				req := api.CalculateRequest{
					Values:     values,
					Percentile: 50 + float64(rand.Intn(50)),
				}

				body, _ := json.Marshal(req)
				httpReq := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
				httpReq.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				srv.router.ServeHTTP(w, httpReq)

				if w.Code != http.StatusOK {
					errors <- fmt.Errorf("goroutine %d, request %d: got status %d", id, j, w.Code)
					continue
				}

				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)
	totalRequests := numGoroutines * requestsPerGoroutine

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Error(err)
		errorCount++
		if errorCount >= 10 {
			t.Log("... (more errors suppressed)")
			break
		}
	}

	if errorCount > 0 {
		t.Fatalf("Failed with %d errors out of %d requests", errorCount, totalRequests)
	}

	t.Logf("Successfully processed %d concurrent requests in %v", successCount, duration)
	t.Logf("Throughput: %.2f requests/sec", float64(totalRequests)/duration.Seconds())
	t.Logf("Average latency: %v per request", duration/time.Duration(totalRequests))
}

// TestStress_LargePayloads tests the server with large data payloads
func TestStress_LargePayloads(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	sizes := []int{1_000, 10_000, 100_000, 1_000_000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			values := make([]float64, size)
			for i := 0; i < size; i++ {
				values[i] = rand.Float64() * 1000
			}

			req := api.CalculateRequest{
				Values:     values,
				Percentile: 95.0,
			}

			body, err := json.Marshal(req)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			start := time.Now()
			httpReq := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
			httpReq.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			srv.router.ServeHTTP(w, httpReq)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
			}

			var resp api.CalculateResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			t.Logf("Size: %d values, P%.0f: %.2f, Time: %v, Payload: %d bytes",
				size, resp.Percentile, resp.Result, duration, len(body))

			// Performance expectations
			if duration > 10*time.Second {
				t.Logf("WARNING: Request took longer than expected: %v", duration)
			}
		})
	}
}

// TestStress_MixedOperations tests various endpoints under load
func TestStress_MixedOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	const (
		numGoroutines          = 30
		operationsPerGoroutine = 10
	)

	var wg sync.WaitGroup
	successCount := make(map[string]int)
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Randomly choose operation
				op := rand.Intn(3)

				switch op {
				case 0: // Health check
					httpReq := httptest.NewRequest("GET", "/health", http.NoBody)
					w := httptest.NewRecorder()
					srv.router.ServeHTTP(w, httpReq)
					if w.Code == http.StatusOK {
						mu.Lock()
						successCount["health"]++
						mu.Unlock()
					}

				case 1: // Small calculation
					values := make([]float64, 100)
					for k := range values {
						values[k] = rand.Float64() * 100
					}
					req := api.CalculateRequest{Values: values, Percentile: 95.0}
					body, _ := json.Marshal(req)
					httpReq := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
					httpReq.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()
					srv.router.ServeHTTP(w, httpReq)
					if w.Code == http.StatusOK {
						mu.Lock()
						successCount["calculate_small"]++
						mu.Unlock()
					}

				case 2: // Large calculation
					values := make([]float64, 10000)
					for k := range values {
						values[k] = rand.Float64() * 1000
					}
					req := api.CalculateRequest{Values: values, Percentile: 99.0}
					body, _ := json.Marshal(req)
					httpReq := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
					httpReq.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()
					srv.router.ServeHTTP(w, httpReq)
					if w.Code == http.StatusOK {
						mu.Lock()
						successCount["calculate_large"]++
						mu.Unlock()
					}
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalOps := numGoroutines * operationsPerGoroutine
	totalSuccess := 0
	for _, count := range successCount {
		totalSuccess += count
	}

	t.Logf("Completed %d mixed operations in %v", totalSuccess, duration)
	t.Logf("Breakdown: health=%d, small=%d, large=%d",
		successCount["health"], successCount["calculate_small"], successCount["calculate_large"])
	t.Logf("Throughput: %.2f ops/sec", float64(totalOps)/duration.Seconds())
}

// BenchmarkHealthEndpoint benchmarks the health endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	httpReq := httptest.NewRequest("GET", "/health", http.NoBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		srv.router.ServeHTTP(w, httpReq)
	}
}

// BenchmarkCalculateEndpoint benchmarks the calculate endpoint
func BenchmarkCalculateEndpoint(b *testing.B) {
	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	sizes := []int{100, 1_000, 10_000}

	for _, size := range sizes {
		values := make([]float64, size)
		for i := 0; i < size; i++ {
			values[i] = rand.Float64() * 1000
		}

		req := api.CalculateRequest{Values: values, Percentile: 95.0}
		body, _ := json.Marshal(req)

		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				httpReq := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
				httpReq.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				srv.router.ServeHTTP(w, httpReq)
			}
		})
	}
}
