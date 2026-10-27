package services

import (
	"backend/internal/api/requests"
	"backend/internal/configs"
	"backend/internal/mocks"
	"backend/internal/models"
	"backend/pkg/apperrors"
	"backend/pkg/security"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestAuthService() (*AuthService, *mocks.MockUserRepository, *mocks.MockRefreshTokenRepository) {
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockRefreshTokenRepository)
	cfg := &configs.Config{
		JWTSecret:               "test-jwt-secret",
		JWTExpiryMinutes:        15,
		RefreshTokenSecret:             "test-refresh-secret",
		RefreshTokenExpiryHours:        168,
		RefreshTokenGracePeriodSeconds: 30,
	}

	service := NewAuthService(userRepo, tokenRepo, cfg)
	return service, userRepo, tokenRepo
}

func TestRegister_Success(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	req := requests.RegisterRequest{
		FullName: "Jane Doe",
		Username: "janedoe",
		Email:    "jane@example.com",
		Password: "Password123!",
	}

	userRepo.On("FindByEmail", ctx, req.Email).Return(nil, apperrors.ErrNotFound)
	userRepo.On("FindByUsername", ctx, req.Username).Return(nil, apperrors.ErrNotFound)
	userRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	tokenRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	res, err := service.Register(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)
	assert.Equal(t, req.Email, res.User.Email)
	assert.Equal(t, req.Username, res.User.Username)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	service, userRepo, _ := newTestAuthService()
	ctx := context.Background()

	req := requests.RegisterRequest{
		FullName: "Jane Doe",
		Username: "janedoe",
		Email:    "existing@example.com",
		Password: "Password123!",
	}

	existingUser := &models.User{ID: uuid.New(), Email: req.Email}
	userRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	res, err := service.Register(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrEmailAlreadyExists)
	assert.Nil(t, res)

	userRepo.AssertExpectations(t)
}

func TestRegister_DuplicateUsername(t *testing.T) {
	service, userRepo, _ := newTestAuthService()
	ctx := context.Background()

	req := requests.RegisterRequest{
		FullName: "Jane Doe",
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "Password123!",
	}

	userRepo.On("FindByEmail", ctx, req.Email).Return(nil, apperrors.ErrNotFound)
	existingUser := &models.User{ID: uuid.New(), Username: req.Username}
	userRepo.On("FindByUsername", ctx, req.Username).Return(existingUser, nil)

	res, err := service.Register(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrUsernameAlreadyExists)
	assert.Nil(t, res)

	userRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	password := "SecretPassword123!"
	hash, _ := security.HashPassword(password)
	user := &models.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		Username:     "validuser",
		PasswordHash: hash,
		IsActive:     true,
	}

	req := requests.LoginRequest{
		Identifier: "user@example.com",
		Password:   password,
	}

	userRepo.On("FindByEmailOrUsername", ctx, req.Identifier).Return(user, nil)
	tokenRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	res, err := service.Login(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	service, userRepo, _ := newTestAuthService()
	ctx := context.Background()

	req := requests.LoginRequest{
		Identifier: "unknown@example.com",
		Password:   "password",
	}

	userRepo.On("FindByEmailOrUsername", ctx, req.Identifier).Return(nil, apperrors.ErrNotFound)

	res, err := service.Login(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
	assert.Nil(t, res)
}

func TestLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	service, userRepo, _ := newTestAuthService()
	ctx := context.Background()

	hash, _ := security.HashPassword("CorrectPassword123!")
	user := &models.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: hash,
		IsActive:     true,
	}

	req := requests.LoginRequest{
		Identifier: "user@example.com",
		Password:   "WrongPassword!",
	}

	userRepo.On("FindByEmailOrUsername", ctx, req.Identifier).Return(user, nil)

	res, err := service.Login(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
	assert.Nil(t, res)
}

func TestLogin_InactiveUser(t *testing.T) {
	service, userRepo, _ := newTestAuthService()
	ctx := context.Background()

	user := &models.User{
		ID:       uuid.New(),
		Email:    "inactive@example.com",
		IsActive: false,
	}

	req := requests.LoginRequest{
		Identifier: "inactive@example.com",
		Password:   "password",
	}

	userRepo.On("FindByEmailOrUsername", ctx, req.Identifier).Return(user, nil)

	res, err := service.Login(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrUserInactive)
	assert.Nil(t, res)
}

func TestRefreshToken_Success(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "sample-refresh-token-string"
	tokenHash := security.HashToken(rawToken)
	userID := uuid.New()

	tokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		IsRevoked: false,
	}

	user := &models.User{
		ID:       userID,
		Email:    "user@example.com",
		Username: "user",
		IsActive: true,
	}

	req := requests.RefreshTokenRequest{RefreshToken: rawToken}

	tokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(tokenRecord, nil)
	userRepo.On("FindByID", ctx, userID.String()).Return(user, nil)
	tokenRepo.On("Revoke", ctx, tokenHash, "rotated").Return(nil)
	tokenRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	res, err := service.RefreshToken(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_GracePeriod_ConcurrentSuccess(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "concurrent-refresh-token"
	tokenHash := security.HashToken(rawToken)
	userID := uuid.New()

	reason := "rotated"
	revokedAt := time.Now().UTC().Add(-5 * time.Second) // revoked 5 seconds ago (well within 30s grace period)
	tokenRecord := &models.RefreshToken{
		ID:               uuid.New(),
		UserID:           userID,
		TokenHash:        tokenHash,
		ExpiresAt:        time.Now().UTC().Add(24 * time.Hour),
		IsRevoked:        true,
		RevokedAt:        &revokedAt,
		RevocationReason: &reason,
	}

	user := &models.User{
		ID:       userID,
		Email:    "concurrent@example.com",
		Username: "concurrent_user",
		IsActive: true,
	}

	req := requests.RefreshTokenRequest{RefreshToken: rawToken}

	tokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(tokenRecord, nil)
	userRepo.On("FindByID", ctx, userID.String()).Return(user, nil)
	tokenRepo.On("Create", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	res, err := service.RefreshToken(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_ReuseDetection_AfterGracePeriod(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "stolen-old-token"
	tokenHash := security.HashToken(rawToken)
	userID := uuid.New()

	reason := "rotated"
	revokedAt := time.Now().UTC().Add(-60 * time.Second) // revoked 60 seconds ago (expired grace period!)
	revokedTokenRecord := &models.RefreshToken{
		ID:               uuid.New(),
		UserID:           userID,
		TokenHash:        tokenHash,
		ExpiresAt:        time.Now().UTC().Add(24 * time.Hour),
		IsRevoked:        true,
		RevokedAt:        &revokedAt,
		RevocationReason: &reason,
	}

	user := &models.User{
		ID:       userID,
		Email:    "victim@example.com",
		Username: "victim_user",
		IsActive: true,
	}

	req := requests.RefreshTokenRequest{RefreshToken: rawToken}

	tokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(revokedTokenRecord, nil)
	userRepo.On("FindByID", ctx, userID.String()).Return(user, nil)
	tokenRepo.On("RevokeAllByUserID", ctx, userID, "reuse_attack_detected").Return(nil)

	res, err := service.RefreshToken(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrTokenRevoked)
	assert.Nil(t, res)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_LoggedOutToken_NoGracePeriod(t *testing.T) {
	service, userRepo, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "logged-out-token"
	tokenHash := security.HashToken(rawToken)
	userID := uuid.New()

	reason := "logged_out"
	revokedAt := time.Now().UTC().Add(-2 * time.Second) // even though 2 seconds ago, it was logged out!
	revokedTokenRecord := &models.RefreshToken{
		ID:               uuid.New(),
		UserID:           userID,
		TokenHash:        tokenHash,
		ExpiresAt:        time.Now().UTC().Add(24 * time.Hour),
		IsRevoked:        true,
		RevokedAt:        &revokedAt,
		RevocationReason: &reason,
	}

	user := &models.User{
		ID:       userID,
		Email:    "loggedout@example.com",
		Username: "loggedout_user",
		IsActive: true,
	}

	req := requests.RefreshTokenRequest{RefreshToken: rawToken}

	tokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(revokedTokenRecord, nil)
	userRepo.On("FindByID", ctx, userID.String()).Return(user, nil)

	res, err := service.RefreshToken(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrTokenRevoked)
	assert.Nil(t, res)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_Expired(t *testing.T) {
	service, _, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "expired-token"
	tokenHash := security.HashToken(rawToken)
	userID := uuid.New()

	expiredTokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour), // Expired
		IsRevoked: false,
	}

	req := requests.RefreshTokenRequest{RefreshToken: rawToken}

	tokenRepo.On("FindByTokenHash", ctx, tokenHash).Return(expiredTokenRecord, nil)

	res, err := service.RefreshToken(ctx, req)
	assert.ErrorIs(t, err, apperrors.ErrTokenExpired)
	assert.Nil(t, res)

	tokenRepo.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	service, _, tokenRepo := newTestAuthService()
	ctx := context.Background()

	rawToken := "logout-token"
	tokenHash := security.HashToken(rawToken)

	req := requests.LogoutRequest{RefreshToken: rawToken}
	tokenRepo.On("Revoke", ctx, tokenHash, "logged_out").Return(nil)

	err := service.Logout(ctx, req)
	assert.NoError(t, err)

	tokenRepo.AssertExpectations(t)
}
