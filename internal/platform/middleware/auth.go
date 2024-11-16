package middleware

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bperezgo/rtsp/internal/platform/repository/auth"
	"github.com/bperezgo/rtsp/shared/domain/observability"
	"go.opentelemetry.io/otel/attribute"
)

const extensionName = "Auth"

type Auth struct {
	AuthRepository auth.InmemoryRepository
	TracerProvider observability.TracerProvider
}

// ExtensionName returns the extension name.
func (a *Auth) ExtensionName() string {
	return extensionName
}

// Validate checks if the extension is configured properly.
func (a *Auth) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (a *Auth) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	newCtx, span := a.TracerProvider.Tracer().Start(ctx, "Auth.Middleware")
	defer span.End()

	if !graphql.HasOperationContext(ctx) {
		return next(ctx)
	}

	user := a.AuthRepository.GetUser()

	span.SetAttributes(attribute.String("user.id", user.ID))
	span.SetAttributes(attribute.String("user.company_id", user.CompanyID))

	newCtx = context.WithValue(newCtx, observability.UserCtxKey, user)

	return next(newCtx)
}
