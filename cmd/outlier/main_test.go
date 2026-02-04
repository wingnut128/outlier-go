package main

import (
	"testing"
)

func TestParseValuesFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []float64
		wantErr bool
	}{
		{
			name:  "valid comma-separated values",
			input: "1.0, 2.0, 3.0",
			want:  []float64{1.0, 2.0, 3.0},
		},
		{
			name:  "no spaces",
			input: "1.0,2.0,3.0",
			want:  []float64{1.0, 2.0, 3.0},
		},
		{
			name:  "mixed spacing",
			input: "1.0,  2.0,   3.0",
			want:  []float64{1.0, 2.0, 3.0},
		},
		{
			name:    "invalid number",
			input:   "1.0, abc, 3.0",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only whitespace",
			input:   "   ",
			wantErr: true,
		},
		{
			name:  "trailing comma",
			input: "1.0, 2.0,",
			want:  []float64{1.0, 2.0},
		},
		{
			name:  "leading comma",
			input: ", 1.0, 2.0",
			want:  []float64{1.0, 2.0},
		},
		{
			name:  "multiple commas",
			input: "1.0,,,2.0",
			want:  []float64{1.0, 2.0},
		},
		{
			name:  "single value",
			input: "42.0",
			want:  []float64{42.0},
		},
		{
			name:  "negative values",
			input: "-1.0, -2.0, -3.0",
			want:  []float64{-1.0, -2.0, -3.0},
		},
		{
			name:  "scientific notation",
			input: "1e2, 2.5e-1, 3.14e0",
			want:  []float64{100.0, 0.25, 3.14},
		},
		{
			name:  "integers",
			input: "1, 2, 3, 4, 5",
			want:  []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		},
		{
			name:    "partial invalid",
			input:   "1.0, 2.0, invalid, 4.0",
			wantErr: true,
		},
		{
			name:  "zero values",
			input: "0, 0.0, 0.00",
			want:  []float64{0.0, 0.0, 0.0},
		},
		{
			name:  "large numbers",
			input: "999999.99, 1000000.00",
			want:  []float64{999999.99, 1000000.00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValuesFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValuesFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !floatSlicesEqual(got, tt.want) {
					t.Errorf("parseValuesFromString() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func floatSlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
