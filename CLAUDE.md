# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Security Constraints

- NEVER read or reference any `.env`, `.env.*`, or files in `secrets/`
- NEVER make network requests, use curl/wget, or fetch external URLs
- ONLY modify files within `./src/` and `./tests/`
- NEVER install new npm packages without explicit approval
- Before any destructive operation (delete, overwrite), state what you plan to do and wait for confirmation

# Working Directories

You have access to:
- `./src/` — application source code (read/write)
- `./tests/` — test files (read/write)
- `./docs/` — documentation (read only)
- `./pkg/` - pkg sources
- `./cmd/` - sources

You do NOT have access to anything outside this project directory.

## Project Overview

Outlier is a percentile calculator written in Go, converted from a Rust implementation. It supports both CLI mode (calculate percentiles from command-line values/files) and server mode (HTTP API for percentile calculations). The core algorithm uses linear interpolation for percentile calculation.

## Build and Test Commands

```bash
# Build
make build                    # Build binary to bin/outlier
make release                  # Build optimized release binary (stripped)

# Testing
make test                     # Run all tests
go test ./internal/calculator -v  # Test specific package
make coverage                 # Generate coverage report (coverage.html)
make stress                   # Run all stress tests
make stress-calc              # Run calculator stress tests only
make stress-server            # Run server stress tests only
make bench                    # Run benchmarks

# Quality
make lint                     # Run golangci-lint (see .golangci.yml for config)

# Development
make clean                    # Remove bin/ and coverage files
make install                  # Install binary to $GOPATH/bin
make hooks                    # Setup git pre-commit/pre-push hooks
```

## Running the Application

```bash
# CLI mode
./bin/outlier --values 1,2,3,4,5 --percentile 50
./bin/outlier --file examples/sample.json --percentile 95

# Server mode
./bin/outlier --serve                                    # Default port 3000
./bin/outlier --serve --port 8080                        # Custom port
./bin/outlier --serve --config configs/config.development.toml
```

## Architecture

### Core Components

**cmd/outlier/main.go**: Entry point using Cobra for CLI. Handles two modes:
- CLI mode: Parses flags, reads input (direct values/files), calculates percentile, outputs result
- Server mode: Initializes HTTP server with Gin framework

**internal/calculator/percentile.go**: Core percentile calculation logic using linear interpolation:
1. Sort values (creates copy to avoid modifying original)
2. Calculate index: `(percentile / 100) * (length - 1)`
3. If exact index, return that value; otherwise interpolate between lower and upper indices

**internal/parser/parser.go**: File parsing for JSON and CSV inputs. CSV files must have a "value" column. Supports both file paths and byte slices (for HTTP uploads).

**internal/server/**: HTTP server implementation:
- `server.go`: Gin router setup, CORS, graceful shutdown, request logging middleware
- `handlers.go`: API endpoints (POST /calculate, POST /calculate/file, GET /health)
- Swagger documentation available at /docs endpoint

**internal/config/config.go**: TOML-based configuration with priority:
1. --config flag
2. CONFIG_FILE environment variable
3. Default configuration (port 3000, info logging)

**internal/telemetry/telemetry.go**: OpenTelemetry integration for Honeycomb (enabled via HONEYCOMB_API_KEY env var)

### Key Design Decisions

- You are a professional Go programmer
- Percentile calculation creates a copy of the input slice before sorting to avoid modifying the original data
- Server uses graceful shutdown with 10-second timeout on SIGTERM/SIGINT
- File uploads limited to 100MB max
- Server timeouts: 10s read header, 30s read/write, 120s idle
- Version info injected at build time via ldflags (see Makefile)

## Testing Conventions

- Test files: `*_test.go` in same package as code under test
- Stress tests: `stress_test.go` files with `TestStress*` naming
- Use table-driven tests for multiple scenarios
- Test naming: `TestFunctionName_Scenario` (e.g., `TestCalculatePercentile_EmptySlice`)
- Target >80% code coverage

## Linter Configuration

See `.golangci.yml` for enabled linters. Key linters:
- errcheck, gosec (security), gocritic, revive, staticcheck
- Excludes G107 (URL in HTTP request) and G404 (weak random in tests)
- Test files exempt from gocyclo, errcheck, gosec
- Complexity threshold: 15

## Commit Conventions

Follow Conventional Commits format:
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Scope: calculator, server, parser, config, etc.
- Update CHANGELOG.md under [Unreleased] section for all changes
- **Do not include "Co-Authored-by" statements in commit messages**

## Common Patterns

**Error Handling**: Return errors with context using `fmt.Errorf("description: %w", err)` for proper error wrapping.

**Configuration**: Use `config.LoadConfigWithPriority()` which respects --config flag > CONFIG_FILE env > defaults.

**File Format Detection**: Uses filepath.Ext() to determine JSON vs CSV, case-insensitive.

**API Response Format**: Consistent JSON structure with count, percentile, result fields (see pkg/api/types.go).


