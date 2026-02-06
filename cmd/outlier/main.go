package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wingnut128/outlier-go/internal/calculator"
	"github.com/wingnut128/outlier-go/internal/config"
	"github.com/wingnut128/outlier-go/internal/parser"
	"github.com/wingnut128/outlier-go/internal/server"
	"github.com/wingnut128/outlier-go/internal/telemetry"
	"github.com/wingnut128/outlier-go/internal/version"
)

// @title Outlier API
// @version 1.0.0
// @description A percentile calculator with CLI and HTTP API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/wingnut128/outlier-go
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /

var (
	serveMode  bool
	configPath string
	port       int
	percentile float64
	filePath   string
	valuesStr  string
)

var rootCmd = &cobra.Command{
	Use:   "outlier",
	Short: "Outlier - Percentile calculator with CLI and HTTP API",
	Long: `Outlier is a percentile calculator that supports both CLI and server modes.
It can calculate percentiles from direct values, JSON files, or CSV files.`,
	Version: version.GetFullVersion(),
	RunE:    runMain,
}

func init() {
	rootCmd.Flags().BoolVar(&serveMode, "serve", false, "Start HTTP API server")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.Flags().IntVar(&port, "port", 0, "Override server port")
	rootCmd.Flags().Float64VarP(&percentile, "percentile", "p", 95.0, "Percentile to calculate (0-100)")
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "Input file path (JSON or CSV)")
	rootCmd.Flags().StringVarP(&valuesStr, "values", "v", "", "Comma-separated values")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runMain(cmd *cobra.Command, args []string) error {
	// Quick exit for help/version without telemetry
	if cmd.Flags().Changed("help") || cmd.Flags().Changed("version") {
		return nil
	}

	// Initialize telemetry
	if err := telemetry.InitTelemetry(); err != nil {
		log.Printf("Warning: failed to initialize telemetry: %v\n", err)
	}
	defer func() {
		if err := telemetry.ShutdownTelemetry(); err != nil {
			log.Printf("Warning: failed to shutdown telemetry: %v\n", err)
		}
	}()

	// Load configuration
	cfg, err := config.LoadConfigWithPriority(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override port if specified
	if port > 0 {
		cfg.Server.Port = port
	}

	// Server mode
	if serveMode {
		return runServer(cfg)
	}

	// CLI mode
	return runCLI()
}

func runServer(cfg *config.Config) error {
	srv := server.NewServer(cfg)
	return srv.Start()
}

func runCLI() error {
	var values []float64
	var err error

	// Determine input source
	switch {
	case filePath != "":
		values, err = parser.ReadValuesFromFile(filePath)
		if err != nil {
			return err
		}
	case valuesStr != "":
		values, err = parseValuesFromString(valuesStr)
		if err != nil {
			return fmt.Errorf("parsing values: %w", err)
		}
	default:
		return fmt.Errorf("must provide either --file or --values")
	}

	// Calculate percentile
	result, err := calculator.CalculatePercentile(values, percentile)
	if err != nil {
		return err
	}

	// Output result
	fmt.Printf("Number of values: %d\n", len(values))
	fmt.Printf("Percentile (P%.0f): %.2f\n", percentile, result)
	return nil
}

func parseValuesFromString(s string) ([]float64, error) {
	parts := strings.Split(s, ",")
	values := make([]float64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		value, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", part)
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no values provided")
	}

	return values, nil
}
