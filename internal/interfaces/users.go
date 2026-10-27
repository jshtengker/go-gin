package interfaces

import (
	"backend/internal/api/responses"
	"backend/internal/models"
	"backend/pkg/helpers"
	"context"
)

type UserRepository interface {
	FindAll(ctx context.Context, pagination helpers.PaginationParams) ([]models.User, int64, error)
	FindByID(ctx context.Context, userID string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type UserService interface {
	GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.UserResponse, int64, error)
	GetByID(ctx context.Context, userID string) (*responses.UserResponse, error)
}
