package observability

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	End()
	IsRecording() bool
	SetAttributes(kv ...attribute.KeyValue)
	SetStatus(code codes.Code, description string)
	RecordError(err error, options ...trace.EventOption)
}

type SpanImpl struct {
	span trace.Span
}

func NewSpan(span trace.Span) Span {
	return &SpanImpl{
		span: span,
	}
}

func (s *SpanImpl) End() {
	s.span.End()
}

func (s *SpanImpl) IsRecording() bool {
	return s.span.IsRecording()
}

func (s *SpanImpl) SetAttributes(kv ...attribute.KeyValue) {
	s.span.SetAttributes(kv...)
}

func (s *SpanImpl) SetStatus(code codes.Code, description string) {
	s.span.SetStatus(code, description)
}

func (s *SpanImpl) RecordError(err error, options ...trace.EventOption) {
	s.span.RecordError(err, options...)
}
