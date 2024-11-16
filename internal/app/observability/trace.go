package observability

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tp trace.TracerProvider
}

func New(tp trace.TracerProvider) Tracer {
	return Tracer{
		tp: tp,
	}
}

func (t *Tracer) StartSpan(ctx context.Context, operationName string) (context.Context, SpanFactory) {
	tracer := t.tp.Tracer(operationName)
	ctx, span := tracer.Start(ctx, operationName)

	return ctx, SpanFactory{
		span: span,
	}
}
