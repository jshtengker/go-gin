package handlers

import (
	"backend/internal/services"
	"backend/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
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
