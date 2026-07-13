package responses

import (
	"backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Sku         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Stock       int       `json:"stock"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToProductResponse(product *models.Products) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Sku:         product.Sku,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ToProductsResponse(products []models.Products) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))

	for _, product := range products {
		responses = append(responses, *ToProductResponse(&product))
	}

	return responses
}
