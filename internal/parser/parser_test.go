package parser

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestReadJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")

	content := `[1.0, 2.0, 3.0, 4.0, 5.0]`
	err := os.WriteFile(jsonFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	values, err := ReadJSONFile(jsonFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
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
	err := os.WriteFile(csvFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	values, err := ReadCSVFile(csvFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}

func TestReadCSVFile_MissingValueColumn(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := `number
1.0
2.0`
	err := os.WriteFile(csvFile, []byte(content), 0o644)
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

	err := os.WriteFile(txtFile, []byte("test"), 0o644)
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
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
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
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}

func TestReadValuesFromFile_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "data.json")
	if err := os.WriteFile(jsonFile, []byte(`[10, 20]`), 0o644); err != nil {
		t.Fatal(err)
	}

	values, err := ReadValuesFromFile(jsonFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(values, []float64{10, 20}) {
		t.Errorf("expected [10 20], got %v", values)
	}
}

func TestReadValuesFromFile_CSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "data.csv")
	if err := os.WriteFile(csvFile, []byte("value\n7\n8\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	values, err := ReadValuesFromFile(csvFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(values, []float64{7, 8}) {
		t.Errorf("expected [7 8], got %v", values)
	}
}

func TestReadJSONFile_NotFound(t *testing.T) {
	_, err := ReadJSONFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestReadJSONFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "bad.json")
	if err := os.WriteFile(f, []byte(`not json`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ReadJSONFile(f)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestReadCSVFile_NotFound(t *testing.T) {
	_, err := ReadCSVFile("/nonexistent/file.csv")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestReadJSONBytes_InvalidJSON(t *testing.T) {
	_, err := ReadJSONBytes([]byte(`{not json}`))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestReadValuesFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		filename string
		want     []float64
		wantErr  bool
	}{
		{
			name:     "JSON bytes",
			data:     []byte(`[1, 2, 3]`),
			filename: "upload.json",
			want:     []float64{1, 2, 3},
		},
		{
			name:     "CSV bytes",
			data:     []byte("value\n4\n5\n6\n"),
			filename: "upload.csv",
			want:     []float64{4, 5, 6},
		},
		{
			name:     "case-insensitive extension",
			data:     []byte(`[9]`),
			filename: "FILE.JSON",
			want:     []float64{9},
		},
		{
			name:     "unsupported format",
			data:     []byte("data"),
			filename: "file.xml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadValuesFromBytes(tt.data, tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestReadCSVBytes_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []float64
		wantErr bool
	}{
		{
			name: "empty values skipped",
			data: "value\n1\n\n3\n",
			want: []float64{1, 3},
		},
		{
			name: "multi-column with value column",
			data: "name,value,extra\na,10,x\nb,20,y\n",
			want: []float64{10, 20},
		},
		{
			name:    "invalid number",
			data:    "value\nabc\n",
			wantErr: true,
		},
		{
			name:    "empty CSV",
			data:    "",
			wantErr: true,
		},
		{
			name:    "header only no value column",
			data:    "id,name\n",
			wantErr: true,
		},
		{
			name: "whitespace in header",
			data: " Value \n42\n",
			want: []float64{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCSVBytes([]byte(tt.data))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
