package middlewares

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

const extensionName = "Auth"

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

type User struct {
	ID        string
	CompanyID string
}

type AuthRepository struct{}

func (a *AuthRepository) GetUser() User {
	return User{
		ID:        "1",
		CompanyID: "3",
	}
}

type Auth struct {
	AuthRepository AuthRepository
}

// ExtensionName returns the extension name.
func (a Auth) ExtensionName() string {
	return extensionName
}

// Validate checks if the extension is configured properly.
func (a Auth) Validate(_ graphql.ExecutableSchema) error {
	return nil
}

func (a Auth) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	if !graphql.HasOperationContext(ctx) {
		return next(ctx)
	}

	// Here is the validation of who is the user, and with the next value set in the context,
	// It is possible to add it for each span in the next layers.

	// After that could happen the authorization validation

	// And before this it needs to be the error handling

	user := a.AuthRepository.GetUser()
	ctx = context.WithValue(ctx, userCtxKey, user)

	return next(ctx)
}
