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

func (s *ProductService) GetById(ctx context.Context, product_id string) (*responses.ProductResponse, error) {
	product, err := s.repo.FindById(ctx, product_id)
	if err != nil {
		return nil, err
	}

	return responses.ToProductResponse(product), nil
}
