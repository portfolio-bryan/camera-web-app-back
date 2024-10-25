package middlewares

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	domainerror "github.com/bperezgo/rtsp/shared/domain/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	// var domainErr *domainerror.Domain

	err := graphql.DefaultErrorPresenter(ctx, e)

	switch e := err.Err.(type) {
	case domainerror.Domain:
		err.Message = e.Error()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("error.code", string(e.Code())))
		span.SetAttributes(attribute.String("error.type", string(e.Type())))
	default:
		err.Message = "Service Unavailable"
	}
	return err
}
