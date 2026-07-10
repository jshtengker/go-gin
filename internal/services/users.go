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

	userResponses := make([]responses.UserResponse, 0, len(users))

	for _, user := range users {
		userResponses = append(userResponses, responses.UserResponse{
			ID:          user.ID,
			FullName:    user.FullName,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			IsVerified:  user.IsVerified,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	return userResponses, nil
}
