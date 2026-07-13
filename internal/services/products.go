package services

import (
	"backend/internal/api/responses"
	"backend/internal/repositories"
	"backend/pkg/helpers"
	"context"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
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
