package calculator

import (
	"math"
	"testing"
)

const epsilon = 0.0001

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestCalculatePercentile_SimpleMedian(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	result, err := CalculatePercentile(values, 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 3.0) {
		t.Errorf("expected 3.0, got %.4f", result)
	}
}

func TestCalculatePercentile_P95(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	result, err := CalculatePercentile(values, 95.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 9.55) {
		t.Errorf("expected 9.55, got %.4f", result)
	}
}

func TestCalculatePercentile_P99(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	result, err := CalculatePercentile(values, 99.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 9.91) {
		t.Errorf("expected 9.91, got %.4f", result)
	}
}

func TestCalculatePercentile_P0(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	result, err := CalculatePercentile(values, 0.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 1.0) {
		t.Errorf("expected 1.0, got %.4f", result)
	}
}

func TestCalculatePercentile_P100(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	result, err := CalculatePercentile(values, 100.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 5.0) {
		t.Errorf("expected 5.0, got %.4f", result)
	}
}

func TestCalculatePercentile_EmptySlice(t *testing.T) {
	values := []float64{}
	_, err := CalculatePercentile(values, 50.0)
	if err == nil {
		t.Error("expected error for empty slice, got nil")
	}
}

func TestCalculatePercentile_SingleValue(t *testing.T) {
	values := []float64{42.0}
	result, err := CalculatePercentile(values, 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 42.0) {
		t.Errorf("expected 42.0, got %.4f", result)
	}
}

func TestCalculatePercentile_UnsortedInput(t *testing.T) {
	values := []float64{5.0, 2.0, 8.0, 1.0, 9.0}
	result, err := CalculatePercentile(values, 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 5.0) {
		t.Errorf("expected 5.0, got %.4f", result)
	}
}

func TestCalculatePercentile_Duplicates(t *testing.T) {
	values := []float64{1.0, 2.0, 2.0, 2.0, 3.0}
	result, err := CalculatePercentile(values, 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(result, 2.0) {
		t.Errorf("expected 2.0, got %.4f", result)
	}
}

func TestCalculatePercentile_LargeDataset(t *testing.T) {
	values := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		values[i] = float64(i + 1)
	}
	result, err := CalculatePercentile(values, 95.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 950.05
	if !almostEqual(result, expected) {
		t.Errorf("expected %.2f, got %.4f", expected, result)
	}
}

func TestCalculatePercentile_OutOfRange_Negative(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0}
	_, err := CalculatePercentile(values, -1.0)
	if err == nil {
		t.Error("expected error for negative percentile, got nil")
	}
}

func TestCalculatePercentile_OutOfRange_Above100(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0}
	_, err := CalculatePercentile(values, 101.0)
	if err == nil {
		t.Error("expected error for percentile > 100, got nil")
	}
}

func TestCalculatePercentile_DoesNotModifyOriginal(t *testing.T) {
	values := []float64{5.0, 2.0, 8.0, 1.0, 9.0}
	original := make([]float64, len(values))
	copy(original, values)

	_, err := CalculatePercentile(values, 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := range values {
		if values[i] != original[i] {
			t.Errorf("original slice was modified at index %d: expected %.2f, got %.2f", i, original[i], values[i])
		}
	}
}
