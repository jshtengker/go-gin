package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

const RequestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Header(RequestIDHeader, requestID)
		c.Set(string(RequestIDKey), requestID)

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if val, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
