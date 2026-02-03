package api

// CalculateRequest represents a request to calculate a percentile
type CalculateRequest struct {
	Values     []float64 `json:"values" binding:"required"`
	Percentile float64   `json:"percentile"`
}

// CalculateResponse represents the result of a percentile calculation
type CalculateResponse struct {
	Count      int     `json:"count"`
	Percentile float64 `json:"percentile"`
	Result     float64 `json:"result"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// ValueRecord represents a CSV record with a value field
type ValueRecord struct {
	Value float64 `csv:"value"`
}
