package services

import (
	"backend/internal/api/responses"
	"backend/internal/repositories"
	"backend/pkg/helpers"
	"context"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.OrderListResponse, int64, error) {
	orders, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil

}

func (s *OrderService) GetById(ctx context.Context, order_id string) (*responses.OrderResponse, error) {
	order, err := s.repo.FindById(ctx, order_id)
	if err != nil {
		return nil, err
	}

	return responses.ToOrderResponse(order), nil
}
