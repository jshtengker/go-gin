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

type ProductHandler struct {
	service interfaces.ProductService
}

func NewProductHandler(service interfaces.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	pagination := responses.ParsePagination(c)

	products, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve products",
		)
		return
	}

	data := responses.BuildPagination(
		pagination,
		total,
		products,
	)

	helpers.SuccessResponse(
		c,
		http.StatusOK,
		data,
		"Products retrieved successfully",
	)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	productID := c.Param("id")

	if _, err := uuid.Parse(productID); err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid product ID",
			nil,
		)
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			helpers.ErrorResponse(
				c,
				http.StatusNotFound,
				"Product not found",
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
		product,
		"Product retrieved successfully",
	)
}
