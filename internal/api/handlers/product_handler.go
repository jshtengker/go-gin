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

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	pagination := helpers.ParsePagination(c)

	products, total, err := h.service.GetAll(c.Request.Context(), pagination)
	if err != nil {
		helpers.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to retrieve products",
		)
		return
	}

	data := helpers.BuildPagination(
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

func (h *ProductHandler) GetById(c *gin.Context) {
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

	product, err := h.service.GetById(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		"Product retrieved Successfully",
	)
}
