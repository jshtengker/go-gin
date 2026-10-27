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

func TestOrderService_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()
	pagination := helpers.NewPaginationParams(1, 10)

	t.Run("Success", func(t *testing.T) {
		orders := []models.OrderListItem{
			{ID: uuid.New(), OrderNumber: "ORD-001", Status: "completed", TotalAmount: 199.99},
			{ID: uuid.New(), OrderNumber: "ORD-002", Status: "pending", TotalAmount: 49.99},
		}
		mockRepo.On("FindAll", ctx, pagination).Return(orders, int64(2), nil).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
		assert.Equal(t, "ORD-001", result[0].OrderNumber)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo.On("FindAll", ctx, pagination).Return(nil, int64(0), errors.New("db error")).Once()

		result, total, err := service.GetAll(ctx, pagination)
		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, result)
	})
}

func TestOrderService_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := NewOrderService(mockRepo)
	ctx := context.Background()
	orderID := uuid.New().String()

	t.Run("Success", func(t *testing.T) {
		order := &models.Orders{
			ID:          uuid.MustParse(orderID),
			OrderNumber: "ORD-12345",
			Status:      "shipped",
			TotalAmount: 350.00,
		}
		mockRepo.On("FindByID", ctx, orderID).Return(order, nil).Once()

		result, err := service.GetByID(ctx, orderID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ORD-12345", result.OrderNumber)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, orderID).Return(nil, apperrors.ErrNotFound).Once()

		result, err := service.GetByID(ctx, orderID)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, result)
	})
}
