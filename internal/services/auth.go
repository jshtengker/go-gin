package services

import (
	"backend/internal/api/requests"
	"backend/internal/api/responses"
	"backend/internal/configs"
	"backend/internal/interfaces"
	"backend/internal/models"
	"backend/pkg/apperrors"
	"backend/pkg/security"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	tokenRepo interfaces.RefreshTokenRepository
	cfg       *configs.Config
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	tokenRepo interfaces.RefreshTokenRepository,
	cfg *configs.Config,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, req requests.RegisterRequest) (*responses.AuthResponse, error) {
	// Check if email already exists
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, apperrors.ErrEmailAlreadyExists
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	// Check if username already exists
	if _, err := s.userRepo.FindByUsername(ctx, req.Username); err == nil {
		return nil, apperrors.ErrUsernameAlreadyExists
	} else if !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	// Hash password
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &models.User{
		ID:           uuid.New(),
		FullName:     req.FullName,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		PhoneNumber:  req.PhoneNumber,
		IsVerified:   false,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req requests.LoginRequest) (*responses.AuthResponse, error) {
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, apperrors.ErrUserInactive
	}

	if !security.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user)
}

func (s *AuthService) RefreshToken(ctx context.Context, req requests.RefreshTokenRequest) (*responses.AuthResponse, error) {
	tokenHash := security.HashToken(req.RefreshToken)

	tokenRecord, err := s.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrInvalidToken
		}
		return nil, err
	}

	// Check expiration
	if time.Now().UTC().After(tokenRecord.ExpiresAt) {
		return nil, apperrors.ErrTokenExpired
	}

	// Verify user is still active
	user, err := s.userRepo.FindByID(ctx, tokenRecord.UserID.String())
	if err != nil {
		return nil, apperrors.ErrInvalidToken
	}

	if !user.IsActive {
		return nil, apperrors.ErrUserInactive
	}

	// Revocation & Grace Period Handling
	if tokenRecord.IsRevoked {
		revocationReason := ""
		if tokenRecord.RevocationReason != nil {
			revocationReason = *tokenRecord.RevocationReason
		}

		gracePeriod := 30 * time.Second
		if s.cfg != nil && s.cfg.RefreshTokenGracePeriodSeconds >= 0 {
			gracePeriod = time.Duration(s.cfg.RefreshTokenGracePeriodSeconds) * time.Second
		}

		var revokedTime time.Time
		if tokenRecord.RevokedAt != nil {
			revokedTime = *tokenRecord.RevokedAt
		} else {
			revokedTime = tokenRecord.UpdatedAt
		}

		// If revoked due to token rotation and within grace period:
		// allow issuing a fresh token pair to handle concurrent frontend requests gracefully.
		if revocationReason == "rotated" && time.Now().UTC().Sub(revokedTime) <= gracePeriod {
			return s.issueTokens(ctx, user)
		}

		// Past grace period or explicit logout:
		// Genuine reuse attack -> revoke all active sessions for this user!
		if revocationReason == "rotated" {
			_ = s.tokenRepo.RevokeAllByUserID(ctx, tokenRecord.UserID, "reuse_attack_detected")
		}
		return nil, apperrors.ErrTokenRevoked
	}

	// Invalidate the old refresh token (single-use rotation)
	if err := s.tokenRepo.Revoke(ctx, tokenHash, "rotated"); err != nil {
		return nil, err
	}

	// Issue fresh token pair
	return s.issueTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, req requests.LogoutRequest) error {
	tokenHash := security.HashToken(req.RefreshToken)
	return s.tokenRepo.Revoke(ctx, tokenHash, "logged_out")
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*responses.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return responses.ToUserResponse(user), nil
}

// issueTokens generates a new Access Token and a new Refresh Token in PostgreSQL.
func (s *AuthService) issueTokens(ctx context.Context, user *models.User) (*responses.AuthResponse, error) {
	accessExpiry := time.Duration(s.cfg.JWTExpiryMinutes) * time.Minute
	accessToken, err := security.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		s.cfg.JWTSecret,
		accessExpiry,
	)
	if err != nil {
		return nil, err
	}

	rawRefreshToken, tokenHash, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	refreshExpiry := time.Duration(s.cfg.RefreshTokenExpiryHours) * time.Hour
	refreshTokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(refreshExpiry),
		IsRevoked: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, err
	}

	return &responses.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiry.Seconds()),
		User:         *responses.ToUserResponse(user),
	}, nil
}
