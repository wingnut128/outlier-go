.PHONY: build release test coverage stress stress-calc stress-server bench install docker-build docker-run swagger clean lint version help

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags for version injection
LDFLAGS := -X github.com/wingnut128/outlier-go/internal/version.Version=$(VERSION) \
           -X github.com/wingnut128/outlier-go/internal/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/wingnut128/outlier-go/internal/version.BuildDate=$(BUILD_DATE)

# Build the binary
build:
	@echo "Building outlier..."
	@go build -ldflags="$(LDFLAGS)" -o bin/outlier ./cmd/outlier

# Build optimized release binary
release:
	@echo "Building release binary..."
	@go build -ldflags="-s -w $(LDFLAGS)" -o bin/outlier ./cmd/outlier

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Generate test coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run stress tests
stress:
	@echo "Running all stress tests..."
	@echo "================================"
	@go test -v -run TestStress ./internal/calculator
	@echo ""
	@echo "================================"
	@go test -v -run TestStress ./internal/server
	@echo ""
	@echo "Stress tests complete!"

# Run calculator stress tests only
stress-calc:
	@echo "Running calculator stress tests..."
	@go test -v -run TestStress ./internal/calculator

# Run server stress tests only
stress-server:
	@echo "Running server stress tests..."
	@go test -v -run TestStress ./internal/server

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./internal/calculator
	@echo ""
	@go test -bench=. -benchmem ./internal/server

# Install the binary
install:
	@echo "Installing outlier..."
	@go install ./cmd/outlier

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t outlier:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run --rm outlier:latest --help

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/outlier/main.go -o docs

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Show version information
version:
	@echo "Version:    $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  release       - Build optimized release binary"
	@echo "  test          - Run tests"
	@echo "  coverage      - Generate test coverage report"
	@echo "  stress        - Run all stress tests"
	@echo "  stress-calc   - Run calculator stress tests only"
	@echo "  stress-server - Run server stress tests only"
	@echo "  bench         - Run benchmarks"
	@echo "  install       - Install the binary"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  version       - Show version information"
	@echo "  help          - Show this help message"
