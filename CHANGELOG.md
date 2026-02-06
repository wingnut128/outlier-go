# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-03

### Added
- Initial Go implementation converted from Rust
- Core percentile calculator with linear interpolation algorithm
- CLI mode with support for:
  - Direct value input via `--values` flag
  - JSON file input via `--file` flag
  - CSV file input via `--file` flag
  - Configurable percentile via `--percentile` flag (default: 95)
- HTTP API server with:
  - `POST /calculate` endpoint for JSON array calculations
  - `POST /calculate/file` endpoint for file uploads
  - `GET /health` endpoint for health checks
  - `GET /docs/*` Swagger UI for interactive API documentation
  - CORS support (allow all origins)
  - 100MB body size limit for large datasets
  - JSON/CSV file format support
- Configuration management:
  - TOML-based configuration files
  - Priority: CLI flag > ENV var > defaults
  - Configurable logging (level, output, format)
  - Configurable server (port, bind IP)
  - Multiple config templates (development, production, minimal)
- OpenTelemetry integration:
  - Honeycomb export support via OTLP
  - Configurable via `HONEYCOMB_API_KEY` environment variable
  - Service name configuration via `OTEL_SERVICE_NAME`
  - Silent initialization when API key not set
- Version management:
  - Build-time version injection from git tags
  - Git commit hash tracking
  - Build timestamp tracking
  - Full version display: `v1.0.0 (abc123d, built 2026-02-03T...)`
- Docker support:
  - Multi-stage Dockerfile for optimized images
  - Build-time version injection
  - Alpine-based runtime image
  - Configuration and example files included
- Comprehensive test suite:
  - 13 calculator unit tests (100% coverage)
  - 7 configuration tests (100% coverage)
  - 6 parser tests (72.8% coverage)
  - Overall: 26 tests, 90.9% coverage
- Stress tests and benchmarks:
  - Large dataset tests (10K to 10M values)
  - Concurrent calculation tests (100 goroutines)
  - Memory usage validation
  - Edge case stress tests
  - API concurrent request tests (50 goroutines, 1000 requests)
  - Large payload tests (up to 1M values)
  - Performance benchmarks for calculator and server
- Makefile with targets:
  - `build` - Build development binary
  - `release` - Build optimized release binary
  - `test` - Run all tests
  - `coverage` - Generate coverage report
  - `stress` - Run all stress tests
  - `stress-calc` - Run calculator stress tests
  - `stress-server` - Run server stress tests
  - `bench` - Run benchmarks
  - `install` - Install binary to GOPATH
  - `docker-build` - Build Docker image
  - `docker-run` - Run Docker container
  - `swagger` - Generate Swagger documentation
  - `lint` - Run linter
  - `clean` - Clean build artifacts
  - `version` - Show version information
  - `help` - Show help message
- Documentation:
  - Comprehensive README.md with usage examples
  - VERIFICATION.md with implementation verification
  - STRESS_TEST.md with performance analysis
  - Swagger/OpenAPI specification
  - Example data files (JSON and CSV)

### Performance
- 10M values: 957ms
- 1M values: 87ms
- 100K values: 11ms
- API throughput: 23,911 requests/sec
- Concurrent: 1000 operations in 107ms
- Memory: ~1.5x input size during calculation
- Allocations: 1 per calculation (single slice allocation)

### Dependencies
- `github.com/spf13/cobra` v1.8.1 - CLI framework
- `github.com/gin-gonic/gin` v1.11.0 - HTTP server
- `github.com/gin-contrib/cors` v1.7.2 - CORS middleware
- `github.com/pelletier/go-toml/v2` v2.2.4 - TOML parsing
- `go.opentelemetry.io/otel` v1.40.0 - OpenTelemetry core
- `go.opentelemetry.io/otel/sdk` v1.40.0 - OpenTelemetry SDK
- `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` v1.40.0 - OTLP exporter
- `github.com/swaggo/swag` v1.16.6 - Swagger generation
- `github.com/swaggo/gin-swagger` v1.6.1 - Gin Swagger integration
- `github.com/swaggo/files` v1.0.1 - Swagger file serving
- `github.com/stretchr/testify` v1.11.1 - Testing utilities

## [Unreleased]

## [1.0.3] - 2026-02-06

### Changed
- Add server handler unit tests, improving coverage from 41% to 78% with tests for health, calculate, file upload endpoints, error paths, and request logger

## [1.0.2] - 2026-02-06

### Changed
- Expand parser test coverage from 70% to 94% with tests for `ReadValuesFromBytes`, error paths, and CSV edge cases

## [1.0.1] - 2026-02-06

### Added
- MIT License file

### Changed
- Deduplicate CSV parsing into shared `readCSVFromReader` helper
- Add `badRequest` helper and `defaultPercentile` constant in server handlers
- Replace custom `floatSlicesEqual` with stdlib `slices.Equal` in tests
- Simplify redundant error wrapping in `runCLI`

### Removed
- Unused `ValueRecord` type from parser

### Fixed
- Pre-commit hook: remove false-positive "unhandled errors" check (already covered by `go vet` and `errcheck` linter)
- Pre-commit hook: make debug print statement check a non-blocking warning instead of interactive prompt
- Pre-push hook: make low coverage warning non-blocking instead of interactive prompt
- Git hooks no longer use `read -p` interactive prompts that fail in non-interactive contexts (CI, piped input)

---

## Version History

- **v1.0.3** (2026-02-06) - Improved server test coverage
- **v1.0.2** (2026-02-06) - Improved parser test coverage
- **v1.0.1** (2026-02-06) - Code cleanup and hook fixes
  - Refactored duplicate code in CSV parsing and server handlers
  - Fixed git hooks to work in non-interactive environments
- **v1.0.0** (2026-02-03) - Initial Go release
  - Complete feature parity with Rust implementation
  - Production-ready with comprehensive testing
  - Performance validated with stress tests
  - Full API documentation via Swagger

---

## Migration from Rust

This is the Go implementation of the original Rust project. Key differences:

### Maintained
- ✅ Identical percentile calculation algorithm (linear interpolation)
- ✅ Same API contract (request/response formats)
- ✅ Same configuration format (TOML)
- ✅ Same CLI interface and flags
- ✅ OpenTelemetry/Honeycomb integration
- ✅ Docker support

### Go-Specific Changes
- **HTTP Framework:** Gin instead of Axum/Actix-web
- **CLI Framework:** Cobra instead of Clap
- **Build System:** Go modules instead of Cargo
- **Feature Flags:** No build-time flags (server always included)
- **Performance:** Similar performance characteristics (O(n log n))

### Improvements in Go Version
- ✅ Simpler deployment (single static binary)
- ✅ No separate feature compilation needed
- ✅ Built-in concurrency with goroutines
- ✅ Extensive stress tests and benchmarks
- ✅ Interactive Swagger UI
- ✅ Build-time version injection from git tags

---

## Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

## Contributing

When contributing, please:
1. Update this CHANGELOG.md under the `[Unreleased]` section
2. Follow the format: Added/Changed/Deprecated/Removed/Fixed/Security
3. Include a brief description of your changes
4. Reference any related issues or PRs

---

*Generated with [Keep a Changelog](https://keepachangelog.com/)*
