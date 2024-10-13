package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MetadataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, uuid.NewString())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
