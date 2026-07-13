package services

import (
	"backend/internal/api/responses"
	"backend/internal/repositories"
	"backend/pkg/helpers"
	"context"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
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

func (s *UserService) GetById(ctx context.Context, user_id string) (*responses.UserResponse, error) {
	user, err := s.repo.FindById(ctx, user_id)
	if err != nil {
		return nil, err
	}

	return responses.ToUserResponse(user), nil

}
