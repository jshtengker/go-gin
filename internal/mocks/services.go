package mocks

import (
	"backend/internal/api/requests"
	"backend/internal/api/responses"
	"backend/pkg/helpers"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req requests.RegisterRequest) (*responses.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req requests.LoginRequest) (*responses.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, req requests.RefreshTokenRequest) (*responses.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, req requests.LogoutRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockAuthService) GetProfile(ctx context.Context, userID string) (*responses.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.UserResponse, int64, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]responses.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) GetByID(ctx context.Context, userID string) (*responses.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.UserResponse), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.ProductResponse, int64, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]responses.ProductResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductService) GetByID(ctx context.Context, productID string) (*responses.ProductResponse, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.ProductResponse), args.Error(1)
}

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetAll(ctx context.Context, pagination helpers.PaginationParams) ([]responses.OrderListResponse, int64, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]responses.OrderListResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderService) GetByID(ctx context.Context, orderID string) (*responses.OrderResponse, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*responses.OrderResponse), args.Error(1)
}
