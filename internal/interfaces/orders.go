package interfaces

import (
	"backend/internal/api/responses"
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"
)

type OrderRepository interface {
	FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]models.OrderListItem, int64, error)
	FindByID(ctx context.Context, orderID string) (*models.Orders, error)
}

type OrderService interface {
	GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.OrderListResponse, int64, error)
	GetByID(ctx context.Context, orderID string) (*responses.OrderResponse, error)
}
