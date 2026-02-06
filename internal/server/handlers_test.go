package server

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wingnut128/outlier-go/internal/config"
	"github.com/wingnut128/outlier-go/pkg/api"
)

func newTestServer() *Server {
	return NewServer(config.DefaultConfig())
}

// --- Health endpoint ---

func TestHandleHealth(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp api.HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %q", resp.Status)
	}
	if resp.Service != "outlier" {
		t.Errorf("expected service 'outlier', got %q", resp.Service)
	}
}

// --- Calculate endpoint ---

func TestHandleCalculate_Success(t *testing.T) {
	srv := newTestServer()
	body := `{"values":[1,2,3,4,5],"percentile":50}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Count != 5 {
		t.Errorf("expected count 5, got %d", resp.Count)
	}
	if resp.Percentile != 50 {
		t.Errorf("expected percentile 50, got %f", resp.Percentile)
	}
	if resp.Result != 3 {
		t.Errorf("expected result 3, got %f", resp.Result)
	}
}

func TestHandleCalculate_DefaultPercentile(t *testing.T) {
	srv := newTestServer()
	body := `{"values":[1,2,3,4,5]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Percentile != 95 {
		t.Errorf("expected default percentile 95, got %f", resp.Percentile)
	}
}

func TestHandleCalculate_InvalidJSON(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp api.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestHandleCalculate_MissingValues(t *testing.T) {
	srv := newTestServer()
	body := `{"percentile":50}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculate_EmptyValues(t *testing.T) {
	srv := newTestServer()
	body := `{"values":[],"percentile":50}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculate_InvalidPercentile(t *testing.T) {
	srv := newTestServer()
	body := `{"values":[1,2,3],"percentile":150}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculate_SingleValue(t *testing.T) {
	srv := newTestServer()
	body := `{"values":[42],"percentile":99}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Result != 42 {
		t.Errorf("expected result 42, got %f", resp.Result)
	}
}

// --- Calculate file endpoint ---

func createMultipartRequest(t *testing.T, filename string, content []byte, percentile string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write file content: %v", err)
	}

	if percentile != "" {
		if err := writer.WriteField("percentile", percentile); err != nil {
			t.Fatalf("failed to write percentile field: %v", err)
		}
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/calculate/file", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestHandleCalculateFile_JSON(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.json", []byte(`[1,2,3,4,5]`), "50")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Count != 5 {
		t.Errorf("expected count 5, got %d", resp.Count)
	}
	if resp.Result != 3 {
		t.Errorf("expected result 3, got %f", resp.Result)
	}
}

func TestHandleCalculateFile_CSV(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.csv", []byte("value\n10\n20\n30\n"), "50")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Count != 3 {
		t.Errorf("expected count 3, got %d", resp.Count)
	}
	if resp.Result != 20 {
		t.Errorf("expected result 20, got %f", resp.Result)
	}
}

func TestHandleCalculateFile_DefaultPercentile(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.json", []byte(`[1,2,3,4,5]`), "")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp api.CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Percentile != 95 {
		t.Errorf("expected default percentile 95, got %f", resp.Percentile)
	}
}

func TestHandleCalculateFile_MissingFile(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/calculate/file", http.NoBody)
	req.Header.Set("Content-Type", "multipart/form-data")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculateFile_UnsupportedFormat(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.xml", []byte("<data/>"), "50")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculateFile_InvalidPercentile(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.json", []byte(`[1,2,3]`), "not_a_number")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculateFile_InvalidJSON(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.json", []byte(`not json`), "50")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandleCalculateFile_EmptyValues(t *testing.T) {
	srv := newTestServer()
	w := httptest.NewRecorder()
	req := createMultipartRequest(t, "data.json", []byte(`[]`), "50")
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

// --- Server setup ---

func TestNewServer_DebugMode(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Logging.Level = "debug"
	srv := NewServer(cfg)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestNewServer_ReleaseMode(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Logging.Level = "info"
	srv := NewServer(cfg)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
}

// --- Request logger ---

func TestRequestLogger_DefaultFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Logging.Format = "text"
	srv := NewServer(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequestLogger_JSONFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Logging.Format = "json"
	srv := NewServer(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequestLogger_WithQueryString(t *testing.T) {
	cfg := config.DefaultConfig()
	srv := NewServer(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health?foo=bar", http.NoBody)
	srv.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
