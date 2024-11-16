package observability

import (
	sharedob "github.com/bperezgo/rtsp/shared/domain/observability"
)

type SpanFactory struct {
	span sharedob.Span
}

func (s *SpanFactory) End() {
	s.span.End()
}

func (s *SpanFactory) AddError(err error) {}

func (s *SpanFactory) WrapError(err error) error {
	return err
}
