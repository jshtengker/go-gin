package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ErrUnauthorizedContext = errors.New("user not found in request context")

func GetAuthUserID(c *gin.Context) (uuid.UUID, error) {
	val, exists := c.Get("user_uuid")
	if !exists {
		return uuid.Nil, ErrUnauthorizedContext
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUnauthorizedContext
	}

	return id, nil
}
