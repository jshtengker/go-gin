package middlewares

import (
	"backend/pkg/security"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-jwt-secret-key"

	userID := uuid.New()
	email := "test@example.com"
	username := "testuser"

	validToken, err := security.GenerateAccessToken(userID, email, username, secret, 15*time.Minute)
	assert.NoError(t, err)

	expiredToken, err := security.GenerateAccessToken(userID, email, username, secret, -1*time.Minute)
	assert.NoError(t, err)

	wrongSecretToken, err := security.GenerateAccessToken(userID, email, username, "wrong-secret", 15*time.Minute)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectContext  bool
	}{
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectContext:  false,
		},
		{
			name:           "Invalid Header Format (No Bearer prefix)",
			authHeader:     "Basic some-token",
			expectedStatus: http.StatusUnauthorized,
			expectContext:  false,
		},
		{
			name:           "Expired Token",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectContext:  false,
		},
		{
			name:           "Invalid Signature (Wrong Secret)",
			authHeader:     "Bearer " + wrongSecretToken,
			expectedStatus: http.StatusUnauthorized,
			expectContext:  false,
		},
		{
			name:           "Valid Bearer Token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectContext:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(Authenticate(secret))

			var capturedUserID string
			var capturedEmail string
			r.GET("/protected", func(c *gin.Context) {
				capturedUserID = c.GetString("user_id")
				capturedEmail = c.GetString("user_email")
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectContext {
				assert.Equal(t, userID.String(), capturedUserID)
				assert.Equal(t, email, capturedEmail)
			}
		})
	}
}
