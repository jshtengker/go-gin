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

func TestProductService_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()
	pagination := helpers.NewPaginationParams(1, 10)

	t.Run("Success", func(t *testing.T) {
		products := []models.Products{
			{ID: uuid.New(), Name: "Keyboard", Price: "100.00"},
			{ID: uuid.New(), Name: "Mouse", Price: "50.00"},
		}
		mockRepo.On("FindAll", ctx, pagination).Return(products, int64(2), nil).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
		assert.Equal(t, "Keyboard", result[0].Name)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo.On("FindAll", ctx, pagination).Return(nil, int64(0), errors.New("db error")).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, result)
	})
}

func TestProductService_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)
	ctx := context.Background()
	productID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		product := &models.Products{
			ID:    uuid.MustParse(productID),
			Name:  "Monitor",
			Price: "250.00",
		}
		mockRepo.On("FindByID", ctx, productID).Return(product, nil).Once()

		result, err := service.GetByID(ctx, productID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Monitor", result.Name)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, productID).Return(nil, apperrors.ErrNotFound).Once()

		result, err := service.GetByID(ctx, productID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, result)
	})
}
