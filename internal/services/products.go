package services

import (
	"backend/internal/api/responses"
	"backend/internal/interfaces"
	"backend/pkg/helpers"
	"context"
)

type ProductService struct {
	repo interfaces.ProductRepository
}

func NewProductService(repo interfaces.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.ProductResponse, int64, error) {
	products, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	return responses.ToProductsResponse(products), total, nil
}

func (s *ProductService) GetByID(ctx context.Context, productID string) (*responses.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return responses.ToProductResponse(product), nil
}
