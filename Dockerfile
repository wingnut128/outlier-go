# Build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version injection
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

# Build the binary with version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
    -X github.com/wingnut128/outlier-go/internal/version.Version=${VERSION} \
    -X github.com/wingnut128/outlier-go/internal/version.GitCommit=${GIT_COMMIT} \
    -X github.com/wingnut128/outlier-go/internal/version.BuildDate=${BUILD_DATE}" \
    -o outlier ./cmd/outlier

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/outlier /usr/local/bin/outlier

# Copy configuration files
COPY --from=builder /app/configs /configs

# Copy examples (optional, for testing)
COPY --from=builder /app/examples /examples

ENTRYPOINT ["outlier"]
CMD ["--help"]
