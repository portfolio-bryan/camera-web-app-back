package observability

import (
	"context"

	sharedob "github.com/bperezgo/rtsp/shared/domain/observability"
)

type Tracer struct {
	tracer sharedob.Tracer
}

func New(tp sharedob.TracerProvider) Tracer {
	tracer := tp.Tracer()

	return Tracer{
		tracer: tracer,
	}
}

func (t *Tracer) StartSpan(ctx context.Context, operationName string) (context.Context, SpanFactory) {
	ctx, span := t.tracer.Start(ctx, operationName)

	return ctx, SpanFactory{
		span: span,
	}
}
