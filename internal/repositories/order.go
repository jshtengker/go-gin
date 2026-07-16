package repositories

import (
	"backend/internal/api/responses"
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.OrderListResponse, int64, error) {
	var orders []responses.OrderListResponse
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.Orders{}).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Table("orders").
		Select(`
		orders.id,
		orders.user_id,
		orders.shipping_address_id,
		orders.order_number,
		orders.status,
		orders.total_amount,
		orders.ordered_at,
		orders.created_at,
		orders.updated_at,
		users.full_name AS user_full_name,
		users.username,
		users.email AS user_email,
		addresses.phone,
		addresses.address_line,
		addresses.city,
		addresses.province,
		addresses.postal_code,
		addresses.country
	`).
		Joins("LEFT JOIN users ON users.id = orders.user_id").
		Joins("LEFT JOIN addresses ON addresses.id = orders.shipping_address_id").
		Limit(pagination.PageSize).
		Offset(pagination.Offset).
		Scan(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) FindById(ctx context.Context, order_id string) (*models.Orders, error) {
	var order models.Orders

	err := r.db.
		WithContext(ctx).
		Where("id = ?", order_id).
		First(&order).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}
