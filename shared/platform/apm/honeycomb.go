package apm

import (
	"context"
	"log"

	"github.com/honeycombio/otel-config-go/otelconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type HoneycombOptions struct {
	Name string
}

type HoneycombTracerProvider struct {
	tracerProvider trace.TracerProvider
	shutdownFn     func()
	opts           HoneycombOptions
}

func NewHoneycombTracerProvider(ctx context.Context, opts HoneycombOptions) *HoneycombTracerProvider {
	if opts.Name == "" {
		log.Fatal("service name is required")
	}

	otelShutdown, err := otelconfig.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	spanExp, err := newSpanExporter(ctx)

	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := newTraceProvider(spanExp, opts)

	otel.SetTracerProvider(tp)

	return &HoneycombTracerProvider{
		tracerProvider: tp,
		shutdownFn:     otelShutdown,
		opts:           opts,
	}
}

func newSpanExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx)
}

// func newMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
// 	returnotlptracegrpc.New(ctx)
// }

func newTraceProvider(exp sdktrace.SpanExporter, opts HoneycombOptions) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opts.Name),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func (p *HoneycombTracerProvider) TracerProvider() trace.TracerProvider {
	return p.tracerProvider
}

func (p *HoneycombTracerProvider) Shutdown(ctx context.Context) error {
	p.shutdownFn()
	return nil
}
