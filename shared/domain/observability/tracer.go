package observability

import (
	"context"

	"github.com/bperezgo/rtsp/shared/domain/dto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type contextUserKey struct {
	name string
}

var UserCtxKey = &contextUserKey{"user"}

type contextTracerIDKey struct {
	name string
}

var XTracerIDCtxKey = &contextTracerIDKey{"xTracerID"}

const (
	XTracerIDHeader = "X-Tracer-Id"
	XTracerIDKey    = attribute.Key("tracer.id")
)

type Tracer interface {
	Start(ctx context.Context, operationName string) (context.Context, Span)
}

type TracerImpl struct {
	tracer trace.Tracer
}

func NewTracer(tracer trace.Tracer) Tracer {
	return &TracerImpl{
		tracer: tracer,
	}
}

func (t *TracerImpl) Start(ctx context.Context, operationName string) (context.Context, Span) {
	attributes := []attribute.KeyValue{
		xTracerIDHeader(ctx),
	}

	user, ok := ctx.Value(UserCtxKey).(dto.User)
	if ok {
		attributes = append(attributes, spanAttributesFromUser(user)...)
	}

	newCtx, span := t.tracer.Start(ctx, operationName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attributes...,
		),
	)

	return newCtx, NewSpan(span)
}

type TracerProvider interface {
	Tracer() Tracer
}

func xTracerIDHeader(ctx context.Context) attribute.KeyValue {
	xTracerID, ok := ctx.Value(XTracerIDCtxKey).(string)
	if !ok {
		return XTracerIDKey.String("")
	}

	return XTracerIDKey.String(xTracerID)
}

func spanAttributesFromUser(user dto.User) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("user.id", user.ID),
		attribute.String("user.company_id", user.CompanyID),
	}
}
