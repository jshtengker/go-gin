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

type OrderHandler struct {
	service interfaces.OrderService
}

func NewOrderHandler(service interfaces.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	pagination := responses.ParsePagination(c)

	orders, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve orders",
		)
		return
	}

	data := responses.BuildPagination(
		pagination,
		total,
		orders,
	)

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		data,
		"Orders retrieved successfully",
	)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	orderID := c.Param("id")

	if _, err := uuid.Parse(orderID); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid order ID",
			nil,
		)
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			helpers.ErrorResponse(
				c,
				http.StatusNotFound,
				"Order not found",
				nil,
			)
			return
		}

		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to get order",
			nil,
		)
		return
	}

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		order,
		"Order retrieved successfully",
	)
}
