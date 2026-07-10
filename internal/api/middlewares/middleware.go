package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !(c.Request.Header.Get("Token") == "auth"){
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "token not found",
			})
			return
		}
		c.Next()
	}
}

