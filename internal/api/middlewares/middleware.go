package middlewares

import (
	"backend/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Token") != "auth" {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Token not found or invalid")
			c.Abort()
			return
		}
		c.Next()
	}
}
