package services

import (
	"backend/internal/api/responses"
	"backend/internal/repositories"
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

func (s *UserService) GetAll(ctx context.Context) ([]responses.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return responses.ToUsersResponse(users), nil
}

func (s *UserService) GetById(ctx context.Context, user_id string) (*responses.UserResponse, error) {
	user, err := s.repo.FindById(ctx, user_id)
	if err != nil {
		return nil, err
	}

	return responses.ToUserResponse(user), nil

}
