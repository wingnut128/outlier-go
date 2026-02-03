package telemetry

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"google.golang.org/grpc/credentials"
)

var tracerProvider *sdktrace.TracerProvider

// InitTelemetry initializes OpenTelemetry tracing with Honeycomb
func InitTelemetry() error {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		// Silently skip telemetry initialization if no API key is set
		return nil
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "outlier"
	}

	ctx := context.Background()

	// Create OTLP exporter with Honeycomb configuration
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("api.honeycomb.io:443"),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		otlptracegrpc.WithHeaders(map[string]string{
			"x-honeycomb-team": apiKey,
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)

	log.Printf("Telemetry initialized for service: %s\n", serviceName)
	return nil
}

// ShutdownTelemetry flushes and shuts down the tracer provider
func ShutdownTelemetry() error {
	if tracerProvider == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tracerProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	return nil
}
