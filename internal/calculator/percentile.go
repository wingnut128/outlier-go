package calculator

import (
	"fmt"
	"sort"
)

// CalculatePercentile calculates the percentile value using linear interpolation.
// The percentile should be between 0 and 100.
// Returns an error if the values slice is empty or percentile is out of range.
func CalculatePercentile(values []float64, percentile float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("cannot calculate percentile of empty dataset")
	}

	if percentile < 0 || percentile > 100 {
		return 0, fmt.Errorf("percentile must be between 0 and 100, got %.2f", percentile)
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate the index position
	index := (percentile / 100.0) * float64(len(sorted)-1)

	// If index is an exact match, return that value
	lowerIndex := int(index)
	if index == float64(lowerIndex) {
		return sorted[lowerIndex], nil
	}

	// Linear interpolation between lower and upper indices
	upperIndex := lowerIndex + 1
	fraction := index - float64(lowerIndex)
	result := sorted[lowerIndex] + (sorted[upperIndex]-sorted[lowerIndex])*fraction

	return result, nil
}
