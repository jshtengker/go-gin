package repositories

import (
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"

	"gorm.io/gorm"
)

// Products

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]models.Products, int64, error) {
	var products []models.Products
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Products{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pagination.PageSize).Offset(pagination.Offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil

}

func (r *ProductRepository) FindById(ctx context.Context, product_id string) (*models.Products, error) {
	var product models.Products

	err := r.db.WithContext(ctx).Where("id = ?", product_id).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
