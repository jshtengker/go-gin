package handlers

import (
	"net/http"

	"backend/internal/services"
	"backend/pkg/helpers"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve users",
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		users,
		"Users retrieved successfully",
	)
}
