package honeycomb

import (
	"context"
	"log"

	"github.com/bperezgo/rtsp/shared/constants"
	"github.com/bperezgo/rtsp/shared/domain/observability"
	"github.com/honeycombio/otel-config-go/otelconfig"
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"

	// "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	// "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Options struct {
	Name string
}

type HoneycombTracerProvider struct {
	tracerProvider trace.TracerProvider
	tracer         trace.Tracer
	shutdownFn     func()
	opts           Options
}

var _ interface {
	observability.TracerProvider
} = (*HoneycombTracerProvider)(nil)

func NewTracerProvider(ctx context.Context, opts Options) *HoneycombTracerProvider {
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

	res := newResource(opts)

	tp := newTraceProvider(res, spanExp, opts)

	otel.SetTracerProvider(tp)

	tracer := tp.Tracer(
		constants.TracerName,
		trace.WithInstrumentationVersion(otelcontrib.Version()),
	)

	return &HoneycombTracerProvider{
		tracerProvider: tp,
		tracer:         tracer,
		shutdownFn:     otelShutdown,
		opts:           opts,
	}
}

func (p *HoneycombTracerProvider) Tracer() observability.Tracer {
	return observability.NewTracer(p.tracer)
}

func newSpanExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx)
}

// func newMetricExporter(ctx context.Context, res *resource.Resource) *metric.MeterProvider {
// 	metricExp, err := otlpmetricgrpc.New(ctx)

// 	if err != nil {
// 		log.Fatalf("failed to create exporter: %v", err)
// 	}

// 	return metric.NewMeterProvider(
// 		metric.WithResource(res),
// 		metric.WithReader(metric.NewPeriodicReader(metricExp)),
// 	)
// }

func newResource(opts Options) *resource.Resource {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opts.Name),
		),
	)

	if err != nil {
		panic(err)
	}

	return res
}

func newTraceProvider(res *resource.Resource, exp sdktrace.SpanExporter, opts Options) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
}

func (p *HoneycombTracerProvider) TracerProvider() trace.TracerProvider {
	return p.tracerProvider
}

func (p *HoneycombTracerProvider) Shutdown(ctx context.Context) error {
	p.shutdownFn()
	return nil
}
