package handlers

import (
	"errors"
	"net/http"

	"backend/internal/services"
	"backend/pkg/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	pagination := helpers.ParsePagination(c)

	users, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve users",
		)
		return
	}

	data := helpers.BuildPagination(
		pagination,
		total,
		users,
	)

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		data,
		"Users retrieved successfully",
	)
}

func (h *UserHandler) GetById(c *gin.Context) {
	userID := c.Param("id")

	if _, err := uuid.Parse(userID); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid user ID",
			nil,
		)
		return
	}

	user, err := h.service.GetById(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.ErrorResponse(
				c,
				http.StatusNotFound,
				"User not found",
				nil,
			)
			return
		}

		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get user",
			nil,
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		user,
		"User retrieved successfully",
	)
}
