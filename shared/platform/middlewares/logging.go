package middlewares

import (
	"bytes"
	"io"

	"github.com/bperezgo/rtsp/shared/platform/handlertypes"
	"github.com/bperezgo/rtsp/shared/platform/logger"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string

		ctx := c.Request.Context()
		ctx.Value(RequestIDKey)

		requestID, ok := ctx.Value(RequestIDKey).(string)
		if !ok {
			requestID = ""
		}

		byteBody, err := io.ReadAll(c.Request.Body)

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": "Error reading request body",
			})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))

		l := logger.GetLogger()
		l.Info(logger.LogInput{
			Action: "REQUEST",
			State:  logger.SUCCESS,
			Http: &logger.LogHttpInput{
				Request: handlertypes.Request{
					Body: string(byteBody),
				},
			},
			Meta: &handlertypes.Meta{
				RequestId: requestID,
			},
		})

		c.Next()

		l.Info(logger.LogInput{
			Action: "RESPONSE",
			State:  logger.SUCCESS,
			Http: &logger.LogHttpInput{
				Request: handlertypes.Request{
					Body: string(byteBody),
				},
			},
			Meta: &handlertypes.Meta{
				RequestId: requestID,
			},
		})

	}
}
