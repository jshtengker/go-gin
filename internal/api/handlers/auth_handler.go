package handlers

import (
	"backend/internal/api/requests"
	"backend/internal/interfaces"
	"backend/pkg/apperrors"
	"backend/pkg/helpers"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service interfaces.AuthService
}

func NewAuthHandler(service interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register handles new user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrEmailAlreadyExists):
			helpers.ErrorResponse(c, http.StatusConflict, "Email is already registered")
		case errors.Is(err, apperrors.ErrUsernameAlreadyExists):
			helpers.ErrorResponse(c, http.StatusConflict, "Username is already taken")
		default:
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	helpers.SuccessResponse(c, http.StatusCreated, res, "User registered successfully")
}

// Login authenticates a user and returns an access & refresh token pair.
func (h *AuthHandler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid email/username or password")
		case errors.Is(err, apperrors.ErrUserInactive):
			helpers.ErrorResponse(c, http.StatusForbidden, "Account is inactive")
		default:
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to authenticate")
		}
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, res, "Login successful")
}

// RefreshToken rotates the refresh token and returns a new token pair.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req requests.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrTokenRevoked):
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Refresh token was revoked or reuse detected. Please log in again.")
		case errors.Is(err, apperrors.ErrTokenExpired):
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Refresh token has expired. Please log in again.")
		case errors.Is(err, apperrors.ErrInvalidToken):
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		case errors.Is(err, apperrors.ErrUserInactive):
			helpers.ErrorResponse(c, http.StatusForbidden, "Account is inactive")
		default:
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token")
		}
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, res, "Token refreshed successfully")
}

// Logout revokes the provided refresh token in PostgreSQL.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req requests.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	if err := h.service.Logout(c.Request.Context(), req); err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, nil, "Logged out successfully")
}

// GetProfile returns the authenticated user's profile.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, err := helpers.GetAuthUserID(c)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	profile, err := h.service.GetProfile(c.Request.Context(), userID.String())
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			helpers.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, profile, "Profile retrieved successfully")
}
