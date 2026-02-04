package parser

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadValuesFromFile reads values from a file based on its extension
func ReadValuesFromFile(path string) ([]float64, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return ReadJSONFile(path)
	case ".csv":
		return ReadCSVFile(path)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (supported: .json, .csv)", ext)
	}
}

// ReadJSONFile reads a JSON file containing an array of numbers
func ReadJSONFile(path string) ([]float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var values []float64
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return values, nil
}

// ReadCSVFile reads a CSV file with a "value" column
func ReadCSVFile(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Find "value" column index
	valueIndex := -1
	for i, col := range header {
		if strings.TrimSpace(strings.ToLower(col)) == "value" {
			valueIndex = i
			break
		}
	}

	if valueIndex == -1 {
		return nil, fmt.Errorf("CSV file must have a 'value' column")
	}

	// Read values
	var values []float64
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		if valueIndex >= len(record) {
			continue
		}

		valueStr := strings.TrimSpace(record[valueIndex])
		if valueStr == "" {
			continue
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number in CSV: %s", valueStr)
		}

		values = append(values, value)
	}

	return values, nil
}

// ReadValuesFromBytes reads values from a byte slice based on the filename extension
func ReadValuesFromBytes(data []byte, filename string) ([]float64, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return ReadJSONBytes(data)
	case ".csv":
		return ReadCSVBytes(data)
	default:
		return nil, fmt.Errorf("unsupported file format: %s (supported: .json, .csv)", ext)
	}
}

// ReadJSONBytes reads JSON data from a byte slice
func ReadJSONBytes(data []byte) ([]float64, error) {
	var values []float64
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return values, nil
}

// ReadCSVBytes reads CSV data from a byte slice
func ReadCSVBytes(data []byte) ([]float64, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Find "value" column index
	valueIndex := -1
	for i, col := range header {
		if strings.TrimSpace(strings.ToLower(col)) == "value" {
			valueIndex = i
			break
		}
	}

	if valueIndex == -1 {
		return nil, fmt.Errorf("CSV file must have a 'value' column")
	}

	// Read values
	var values []float64
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		if valueIndex >= len(record) {
			continue
		}

		valueStr := strings.TrimSpace(record[valueIndex])
		if valueStr == "" {
			continue
		}

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number in CSV: %s", valueStr)
		}

		values = append(values, value)
	}

	return values, nil
}
