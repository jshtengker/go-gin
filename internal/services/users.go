package services

import (
	"backend/internal/api/responses"
	"backend/internal/interfaces"
	"backend/pkg/helpers"
	"context"
)

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.UserResponse, int64, error) {
	users, total, err := s.repo.FindAll(ctx, pagination)
	if err != nil {
		return nil, 0, err
	}

	return responses.ToUsersResponse(users), total, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*responses.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return responses.ToUserResponse(user), nil
}
