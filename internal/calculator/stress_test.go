package calculator

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TestStress_LargeDataset tests performance with large datasets
func TestStress_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	sizes := []int{10_000, 100_000, 1_000_000, 10_000_000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// Generate random data
			values := make([]float64, size)
			for i := 0; i < size; i++ {
				values[i] = rand.Float64() * 1000
			}

			start := time.Now()
			result, err := CalculatePercentile(values, 95.0)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("Size: %d values, P95: %.2f, Time: %v", size, result, duration)

			// Performance expectations
			if duration > 5*time.Second {
				t.Logf("WARNING: Calculation took longer than expected: %v", duration)
			}
		})
	}
}

// TestStress_ConcurrentCalculations tests thread safety with concurrent operations
func TestStress_ConcurrentCalculations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	const (
		numGoroutines = 100
		dataSize      = 10_000
		iterations    = 10
	)

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*iterations)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine generates its own data
			values := make([]float64, dataSize)
			for j := 0; j < dataSize; j++ {
				values[j] = rand.Float64() * 1000
			}

			// Perform multiple calculations
			for iter := 0; iter < iterations; iter++ {
				percentile := 50.0 + float64(iter)*5.0
				_, err := CalculatePercentile(values, percentile)
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, iteration %d: %v", id, iter, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)
	totalCalculations := numGoroutines * iterations

	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Error(err)
		errorCount++
	}

	if errorCount > 0 {
		t.Fatalf("Failed with %d errors", errorCount)
	}

	t.Logf("Completed %d concurrent calculations in %v", totalCalculations, duration)
	t.Logf("Average: %v per calculation", duration/time.Duration(totalCalculations))
}

// TestStress_MemoryUsage tests memory efficiency
func TestStress_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	const size = 1_000_000

	values := make([]float64, size)
	for i := 0; i < size; i++ {
		values[i] = float64(i)
	}

	// Verify original slice is not modified
	original := make([]float64, len(values))
	copy(original, values)

	_, err := CalculatePercentile(values, 95.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that original is unchanged
	for i := 0; i < len(values); i++ {
		if values[i] != original[i] {
			t.Errorf("original slice was modified at index %d", i)
			break
		}
	}

	t.Logf("Memory test passed: %d values processed without modifying original", size)
}

// TestStress_EdgeCases tests edge cases under stress
func TestStress_EdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	testCases := []struct {
		name       string
		generator  func(size int) []float64
		percentile float64
	}{
		{
			name: "AllZeros",
			generator: func(size int) []float64 {
				return make([]float64, size)
			},
			percentile: 95.0,
		},
		{
			name: "AllOnes",
			generator: func(size int) []float64 {
				values := make([]float64, size)
				for i := range values {
					values[i] = 1.0
				}
				return values
			},
			percentile: 50.0,
		},
		{
			name: "AlternatingValues",
			generator: func(size int) []float64 {
				values := make([]float64, size)
				for i := range values {
					if i%2 == 0 {
						values[i] = 0.0
					} else {
						values[i] = 1.0
					}
				}
				return values
			},
			percentile: 95.0,
		},
		{
			name: "NegativeValues",
			generator: func(size int) []float64 {
				values := make([]float64, size)
				for i := range values {
					values[i] = -rand.Float64() * 1000
				}
				return values
			},
			percentile: 99.0,
		},
		{
			name: "VerySmallRange",
			generator: func(size int) []float64 {
				values := make([]float64, size)
				for i := range values {
					values[i] = 0.0001 + rand.Float64()*0.0001
				}
				return values
			},
			percentile: 75.0,
		},
	}

	const size = 100_000

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			values := tc.generator(size)
			start := time.Now()
			result, err := CalculatePercentile(values, tc.percentile)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			t.Logf("%s: P%.0f = %.6f, Time: %v", tc.name, tc.percentile, result, duration)
		})
	}
}

// BenchmarkCalculatePercentile benchmarks the core function
func BenchmarkCalculatePercentile(b *testing.B) {
	sizes := []int{100, 1_000, 10_000, 100_000}

	for _, size := range sizes {
		values := make([]float64, size)
		for i := 0; i < size; i++ {
			values[i] = rand.Float64() * 1000
		}

		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = CalculatePercentile(values, 95.0)
			}
		})
	}
}

// BenchmarkCalculatePercentile_Parallel benchmarks parallel execution
func BenchmarkCalculatePercentile_Parallel(b *testing.B) {
	const size = 10_000
	values := make([]float64, size)
	for i := 0; i < size; i++ {
		values[i] = rand.Float64() * 1000
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = CalculatePercentile(values, 95.0)
		}
	})
}
