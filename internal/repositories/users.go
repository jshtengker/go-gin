package repositories

import (
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// functionally the same, diff approach - the 'err' could be used elsewhere later
	// err := query.Count(&total).Error
	// if err != nil {
	// 	return nil, 0, err
	// }

	if err := query.Limit(pagination.PageSize).Offset(pagination.Offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil

}

func (r *UserRepository) FindById(ctx context.Context, user_id string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}
