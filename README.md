# Outlier - Percentile Calculator

A high-performance percentile calculator written in Go, supporting both CLI and HTTP API modes. Convert from the original Rust implementation.

## Features

- **Core percentile calculation** with linear interpolation
- **CLI mode** with support for:
  - Direct value input (comma-separated)
  - JSON file input
  - CSV file input
- **HTTP API server** with:
  - RESTful endpoints for percentile calculation
  - File upload support (JSON/CSV)
  - Health check endpoint
  - CORS enabled
- **Configuration management** via TOML files
- **OpenTelemetry integration** for Honeycomb observability
- **Docker support** with multi-stage builds

## Installation

### From Source

```bash
go install github.com/wingnut128/outlier-go/cmd/outlier@latest
```

### Build Locally

```bash
git clone https://github.com/wingnut128/outlier-go
cd outlier-go
make build
```

The binary will be available at `bin/outlier`.

### Docker

```bash
docker build -t outlier:latest .
```

## Usage

### CLI Mode

#### Calculate from direct values

```bash
outlier --values 1,2,3,4,5 --percentile 50
```

Output:
```
Number of values: 5
Percentile (P50): 3.00
```

#### Calculate from JSON file

```bash
outlier --file examples/sample.json --percentile 95
```

#### Calculate from CSV file

```bash
outlier --file examples/sample.csv --percentile 99
```

The CSV file must have a `value` column:
```csv
value
1.0
2.0
3.0
```

### Server Mode

Start the HTTP API server:

```bash
outlier --serve
```

Or with custom configuration:

```bash
outlier --serve --config configs/config.development.toml
```

Or override the port:

```bash
outlier --serve --port 8080
```

### API Endpoints

#### POST /calculate

Calculate percentile from an array of values.

**Request:**
```bash
curl -X POST http://localhost:3000/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "values": [1, 2, 3, 4, 5],
    "percentile": 50
  }'
```

**Response:**
```json
{
  "count": 5,
  "percentile": 50,
  "result": 3
}
```

#### POST /calculate/file

Upload a file (JSON or CSV) and calculate percentile.

**Request:**
```bash
curl -X POST http://localhost:3000/calculate/file \
  -F "file=@examples/sample.json" \
  -F "percentile=95"
```

**Response:**
```json
{
  "count": 100,
  "percentile": 95,
  "result": 95.05
}
```

#### GET /health

Health check endpoint.

**Request:**
```bash
curl http://localhost:3000/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "outlier",
  "version": "1.0.0"
}
```

## Configuration

Configuration is managed via TOML files. Priority order:

1. `--config` flag
2. `CONFIG_FILE` environment variable
3. Default configuration

### Example Configuration

```toml
[logging]
level = "info"      # trace, debug, info, warn, error
output = "stdout"   # stdout, stderr, file
format = "compact"  # compact, pretty, json

[server]
port = 3000
bind_ip = "0.0.0.0"
```

See the `configs/` directory for example configurations:
- `config.example.toml` - Full template
- `config.development.toml` - Development settings
- `config.production.toml` - Production settings
- `config.minimal.toml` - Minimal settings

## OpenTelemetry / Honeycomb Integration

Enable telemetry by setting environment variables:

```bash
export HONEYCOMB_API_KEY=your_api_key_here
export OTEL_SERVICE_NAME=outlier
outlier --serve
```

Traces will be sent to Honeycomb for observability.

## Docker Usage

### Build the image

```bash
make docker-build
```

### Run CLI mode

```bash
docker run --rm outlier:latest --values 1,2,3,4,5 --percentile 50
```

### Run server mode

```bash
docker run --rm -p 3000:3000 outlier:latest --serve
```

### With environment variables

```bash
docker run --rm -p 3000:3000 \
  -e HONEYCOMB_API_KEY=your_key \
  -e OTEL_SERVICE_NAME=outlier-docker \
  outlier:latest --serve
```

## Development

### Run tests

```bash
make test
```

### Generate coverage report

```bash
make coverage
```

This generates `coverage.html` for viewing test coverage.

### Run linter

```bash
make lint
```

Requires `golangci-lint` to be installed.

### Project Structure

```
outlier-go/
├── cmd/outlier/           # CLI entrypoint
├── internal/              # Private application code
│   ├── calculator/        # Percentile calculation logic
│   ├── parser/            # File parsing (JSON/CSV)
│   ├── server/            # HTTP server and handlers
│   ├── config/            # Configuration management
│   └── telemetry/         # OpenTelemetry setup
├── pkg/api/               # Public API types
├── examples/              # Example data files
├── configs/               # Configuration templates
└── docs/                  # Documentation
```

## Algorithm

The percentile calculation uses **linear interpolation**:

1. Sort the input values
2. Calculate index position: `index = (percentile / 100) * (length - 1)`
3. If index is exact, return that value
4. Otherwise, interpolate between the lower and upper indices

This matches the behavior of many statistical packages and provides smooth, accurate results.

## Performance

The implementation is optimized for:
- Large datasets (tested with 1M+ values)
- Low memory footprint (sorts in-place copy)
- Fast HTTP response times
- Concurrent request handling

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass (`make test`)
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Original Implementation

This is a Go port of the original Rust implementation. The core algorithm and API remain identical for compatibility.
