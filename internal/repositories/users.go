package repositories

import (
	"backend/internal/models"
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

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindById(ctx context.Context, user_id string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil

}
