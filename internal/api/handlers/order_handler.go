package handlers

import (
	"backend/internal/services"
	"backend/pkg/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	pagination := helpers.ParsePagination(c)

	orders, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve orders",
		)
		return
	}

	data := helpers.BuildPagination(
		pagination,
		total,
		orders,
	)

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		data,
		"Orders retrieved successfuly",
	)
}

func (h *OrderHandler) GetById(c *gin.Context) {
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

	order, err := h.service.GetById(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
			"Failed to get product",
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
