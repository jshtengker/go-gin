package handlers

import (
	"backend/internal/api/responses"
	"backend/internal/interfaces"
	"backend/pkg/apperrors"
	"backend/pkg/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service interfaces.UserService
}

func NewUserHandler(service interfaces.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	pagination := responses.ParsePagination(c)

	users, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve users",
		)
		return
	}

	data := responses.BuildPagination(
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

func (h *UserHandler) GetByID(c *gin.Context) {
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

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
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
