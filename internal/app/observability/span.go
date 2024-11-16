package observability

import "go.opentelemetry.io/otel/trace"

type SpanFactory struct {
	span trace.Span
}

func (s *SpanFactory) End() {
	s.span.End()
}

func (s *SpanFactory) AddError(err error) {}

func (s *SpanFactory) WrapError(err error) error {
	return err
}
