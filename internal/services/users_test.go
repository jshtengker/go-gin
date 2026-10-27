package services

import (
	"backend/internal/mocks"
	"backend/internal/models"
	"backend/pkg/apperrors"
	"backend/pkg/helpers"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	pagination := helpers.NewPaginationParams(1, 10)

	t.Run("Success", func(t *testing.T) {
		users := []models.User{
			{ID: uuid.New(), FullName: "User 1", Email: "user1@example.com"},
			{ID: uuid.New(), FullName: "User 2", Email: "user2@example.com"},
		}
		mockRepo.On("FindAll", ctx, pagination).Return(users, int64(2), nil).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
		assert.Equal(t, "User 1", result[0].FullName)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("FindAll", ctx, pagination).Return(nil, int64(0), errors.New("db error")).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, result)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	userID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		expectedUser := &models.User{
			ID:       uuid.MustParse(userID),
			FullName: "Found User",
			Email:    "found@example.com",
		}
		mockRepo.On("FindByID", ctx, userID).Return(expectedUser, nil).Once()

		result, err := service.GetByID(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Found User", result.FullName)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, userID).Return(nil, apperrors.ErrNotFound).Once()

		result, err := service.GetByID(ctx, userID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, result)
	})
}
