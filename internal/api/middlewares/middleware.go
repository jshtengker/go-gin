package middlewares

import (
	"backend/pkg/helpers"
	"backend/pkg/security"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticate validates the JWT Bearer access token from the Authorization header.
func Authenticate(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format (expected 'Bearer <token>')")
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		claims, err := security.ValidateAccessToken(tokenStr, jwtSecret)
		if err != nil {
			if errors.Is(err, security.ErrTokenExpired) {
				helpers.ErrorResponse(c, http.StatusUnauthorized, "access token has expired")
			} else {
				helpers.ErrorResponse(c, http.StatusUnauthorized, "invalid or malformed access token")
			}
			c.Abort()
			return
		}

		// Inject authenticated claims into Gin Context
		c.Set("user_id", claims.UserID.String())
		c.Set("user_uuid", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("username", claims.Username)

		c.Next()
	}
}
