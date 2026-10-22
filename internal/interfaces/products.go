package interfaces

import (
	"backend/internal/api/responses"
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"
)

type ProductRepository interface {
	FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]models.Products, int64, error)
	FindByID(ctx context.Context, productID string) (*models.Products, error)
}

type ProductService interface {
	GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.ProductResponse, int64, error)
	GetByID(ctx context.Context, productID string) (*responses.ProductResponse, error)
}
