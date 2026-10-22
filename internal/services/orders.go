package services

import (
	"backend/internal/api/responses"
	"backend/internal/interfaces"
	"backend/pkg/helpers"
	"context"
)

type OrderService struct {
	repo interfaces.OrderRepository
}

func NewOrderService(repo interfaces.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.OrderListResponse, int64, error) {
	items, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	return responses.ToOrderListResponses(items), total, nil
}

func (s *OrderService) GetByID(ctx context.Context, orderID string) (*responses.OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return responses.ToOrderResponse(order), nil
}
