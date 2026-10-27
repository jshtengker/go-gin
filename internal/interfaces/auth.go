package interfaces

import (
	"backend/internal/api/requests"
	"backend/internal/api/responses"
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string, reason string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error
}

type AuthService interface {
	Register(ctx context.Context, req requests.RegisterRequest) (*responses.AuthResponse, error)
	Login(ctx context.Context, req requests.LoginRequest) (*responses.AuthResponse, error)
	RefreshToken(ctx context.Context, req requests.RefreshTokenRequest) (*responses.AuthResponse, error)
	Logout(ctx context.Context, req requests.LogoutRequest) error
	GetProfile(ctx context.Context, userID string) (*responses.UserResponse, error)
}
