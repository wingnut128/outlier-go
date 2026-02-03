package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")

	content := `[1.0, 2.0, 3.0, 4.0, 5.0]`
	err := os.WriteFile(jsonFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	values, err := ReadJSONFile(jsonFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Errorf("at index %d: expected %.2f, got %.2f", i, expected[i], values[i])
		}
	}
}

func TestReadCSVFile(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := `value
1.0
2.0
3.0
4.0
5.0`
	err := os.WriteFile(csvFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	values, err := ReadCSVFile(csvFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Errorf("at index %d: expected %.2f, got %.2f", i, expected[i], values[i])
		}
	}
}

func TestReadCSVFile_MissingValueColumn(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := `number
1.0
2.0`
	err := os.WriteFile(csvFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ReadCSVFile(csvFile)
	if err == nil {
		t.Error("expected error for missing 'value' column, got nil")
	}
}

func TestReadValuesFromFile_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	txtFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(txtFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ReadValuesFromFile(txtFile)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestReadJSONBytes(t *testing.T) {
	data := []byte(`[1.0, 2.0, 3.0]`)
	values, err := ReadJSONBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Errorf("at index %d: expected %.2f, got %.2f", i, expected[i], values[i])
		}
	}
}

func TestReadCSVBytes(t *testing.T) {
	data := []byte(`value
1.0
2.0
3.0`)
	values, err := ReadCSVBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Errorf("at index %d: expected %.2f, got %.2f", i, expected[i], values[i])
		}
	}
}
