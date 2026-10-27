package handlers

import (
	"backend/internal/api/requests"
	"backend/internal/api/responses"
	"backend/internal/mocks"
	"backend/pkg/apperrors"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestAuthRouter(mockService *mocks.MockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewAuthHandler(mockService)

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.RefreshToken)
	r.POST("/logout", handler.Logout)
	r.GET("/me", func(c *gin.Context) {
		// simulate authenticated context if header is provided in test
		if userID := c.GetHeader("X-Test-User-UUID"); userID != "" {
			parsed, _ := uuid.Parse(userID)
			c.Set("user_uuid", parsed)
			c.Set("user_id", parsed.String())
		}
		handler.GetProfile(c)
	})

	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.RegisterRequest{
		FullName: "Alice Smith",
		Username: "alicesmith",
		Email:    "alice@example.com",
		Password: "password123!",
	}

	expectedResponse := &responses.AuthResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}

	mockService.On("Register", mock.Anything, reqPayload).Return(expectedResponse, nil)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	// Password too short (<8 chars)
	invalidPayload := map[string]string{
		"full_name": "Alice",
		"username":  "alice",
		"email":     "not-an-email",
		"password":  "short",
	}

	body, _ := json.Marshal(invalidPayload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Register_ConflictEmail(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.RegisterRequest{
		FullName: "Alice Smith",
		Username: "alicesmith",
		Email:    "existing@example.com",
		Password: "password123!",
	}

	mockService.On("Register", mock.Anything, reqPayload).Return(nil, apperrors.ErrEmailAlreadyExists)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.LoginRequest{
		Identifier: "alice@example.com",
		Password:   "password123!",
	}

	expectedResponse := &responses.AuthResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
	}

	mockService.On("Login", mock.Anything, reqPayload).Return(expectedResponse, nil)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.LoginRequest{
		Identifier: "alice@example.com",
		Password:   "wrongpassword",
	}

	mockService.On("Login", mock.Anything, reqPayload).Return(nil, apperrors.ErrInvalidCredentials)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}

	expectedResponse := &responses.AuthResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
	}

	mockService.On("RefreshToken", mock.Anything, reqPayload).Return(expectedResponse, nil)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_RefreshToken_Revoked(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.RefreshTokenRequest{
		RefreshToken: "revoked-refresh-token",
	}

	mockService.On("RefreshToken", mock.Anything, reqPayload).Return(nil, apperrors.ErrTokenRevoked)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	reqPayload := requests.LogoutRequest{
		RefreshToken: "some-refresh-token",
	}

	mockService.On("Logout", mock.Anything, reqPayload).Return(nil)

	body, _ := json.Marshal(reqPayload)
	req, _ := http.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetProfile_Success(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	userID := uuid.New()
	expectedProfile := &responses.UserResponse{
		ID:       userID,
		Username: "alice",
		Email:    "alice@example.com",
	}

	mockService.On("GetProfile", mock.Anything, userID.String()).Return(expectedProfile, nil)

	req, _ := http.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("X-Test-User-UUID", userID.String())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetProfile_UnauthorizedContext(t *testing.T) {
	mockService := new(mocks.MockAuthService)
	router := setupTestAuthRouter(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/me", nil) // No X-Test-User-UUID header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
