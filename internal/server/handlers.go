package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wingnut128/outlier-go/internal/calculator"
	"github.com/wingnut128/outlier-go/internal/parser"
	"github.com/wingnut128/outlier-go/internal/version"
	"github.com/wingnut128/outlier-go/pkg/api"
)

const defaultPercentile = 95.0

func badRequest(c *gin.Context, format string, args ...any) {
	c.JSON(http.StatusBadRequest, api.ErrorResponse{
		Error: fmt.Sprintf(format, args...),
	})
}

// handleHealth handles GET /health
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} api.HealthResponse
// @Router /health [get]
func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, api.HealthResponse{
		Status:  "healthy",
		Service: "outlier",
		Version: version.GetVersion(),
	})
}

// handleCalculate handles POST /calculate
// @Summary Calculate percentile from values
// @Description Calculate a percentile from an array of numeric values using linear interpolation
// @Tags calculate
// @Accept json
// @Produce json
// @Param request body api.CalculateRequest true "Calculate Request"
// @Success 200 {object} api.CalculateResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /calculate [post]
func handleCalculate(c *gin.Context) {
	var req api.CalculateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, "Invalid request: %v", err)
		return
	}

	// Default percentile to 95 if not provided
	if req.Percentile == 0 {
		req.Percentile = defaultPercentile
	}

	// Calculate percentile
	result, err := calculator.CalculatePercentile(req.Values, req.Percentile)
	if err != nil {
		badRequest(c, "%s", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.CalculateResponse{
		Count:      len(req.Values),
		Percentile: req.Percentile,
		Result:     result,
	})
}

// handleCalculateFile handles POST /calculate/file
// @Summary Calculate percentile from file
// @Description Upload a JSON or CSV file and calculate percentile
// @Tags calculate
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Data file (JSON or CSV)"
// @Param percentile formData number false "Percentile to calculate (default: 95)"
// @Success 200 {object} api.CalculateResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /calculate/file [post]
func handleCalculateFile(c *gin.Context) {
	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		badRequest(c, "Failed to read file: %v", err)
		return
	}
	defer file.Close()

	// Read file contents
	data, err := io.ReadAll(file)
	if err != nil {
		badRequest(c, "Failed to read file contents: %v", err)
		return
	}

	// Parse values from file
	values, err := parser.ReadValuesFromBytes(data, header.Filename)
	if err != nil {
		badRequest(c, "Failed to parse file: %v", err)
		return
	}

	// Get percentile from form or default to 95
	percentile := defaultPercentile
	if percentileStr := c.PostForm("percentile"); percentileStr != "" {
		var p float64
		p, err = strconv.ParseFloat(percentileStr, 64)
		if err != nil {
			badRequest(c, "Invalid percentile value: %v", err)
			return
		}
		percentile = p
	}

	// Calculate percentile
	result, err := calculator.CalculatePercentile(values, percentile)
	if err != nil {
		badRequest(c, "%s", err.Error())
		return
	}

	c.JSON(http.StatusOK, api.CalculateResponse{
		Count:      len(values),
		Percentile: percentile,
		Result:     result,
	})
}
