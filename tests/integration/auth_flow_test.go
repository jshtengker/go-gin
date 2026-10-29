//go:build integration

package integration

import (
	"backend/internal/api/requests"
	"backend/internal/models"
	"backend/pkg/helpers"
	"backend/pkg/security"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apiResponse struct {
	Error   bool                   `json:"error"`
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func TestAuthFlow_Integration(t *testing.T) {
	app, db, _ := setupIntegrationApp(t)

	uniqueID := time.Now().UnixNano()
	testEmail := fmt.Sprintf("testuser_%d@integration.local", uniqueID)
	testUsername := fmt.Sprintf("user_%d", uniqueID)
	testPassword := "SuperSecret123!"

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE email = ?", testEmail)
	})

	var loginAccessToken string
	var loginRefreshToken string
	var rotatedAccessToken string
	var rotatedRefreshToken string

	t.Run("1_Register", func(t *testing.T) {
		regReq := requests.RegisterRequest{
			FullName: "Integration Tester",
			Username: testUsername,
			Email:    testEmail,
			Password: testPassword,
		}
		body, _ := json.Marshal(regReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		require.Equal(t, http.StatusCreated, w.Code, "expected 201 Created on registration: %s", w.Body.String())

		var regRes apiResponse
		err := json.Unmarshal(w.Body.Bytes(), &regRes)
		require.NoError(t, err)

		regAccessToken, ok := regRes.Data["access_token"].(string)
		require.True(t, ok && regAccessToken != "", "expected access token in register response")
		regRefreshToken, ok := regRes.Data["refresh_token"].(string)
		require.True(t, ok && regRefreshToken != "", "expected refresh token in register response")

		var dbUser models.User
		err = db.Where("email = ?", testEmail).First(&dbUser).Error
		require.NoError(t, err, "user should exist in postgres")
		assert.Equal(t, testUsername, dbUser.Username)
	})

	t.Run("2_Login", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Identifier: testEmail,
			Password:   testPassword,
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		require.Equal(t, http.StatusOK, w.Code, "expected 200 OK on login")

		var loginRes apiResponse
		err := json.Unmarshal(w.Body.Bytes(), &loginRes)
		require.NoError(t, err)

		loginAccessToken = loginRes.Data["access_token"].(string)
		loginRefreshToken = loginRes.Data["refresh_token"].(string)
		require.NotEmpty(t, loginAccessToken)
		require.NotEmpty(t, loginRefreshToken)
	})

	t.Run("3_GetProfile_Me", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+loginAccessToken)

		w := executeRequest(app, req)
		require.Equal(t, http.StatusOK, w.Code, "expected 200 OK on /me")

		var profileRes helpers.Response
		err := json.Unmarshal(w.Body.Bytes(), &profileRes)
		require.NoError(t, err)
		profileData := profileRes.Data.(map[string]interface{})
		assert.Equal(t, testEmail, profileData["email"])
		assert.Equal(t, testUsername, profileData["username"])
	})

	t.Run("4_TokenRotation_Refresh", func(t *testing.T) {
		refreshReq := requests.RefreshTokenRequest{
			RefreshToken: loginRefreshToken,
		}
		body, _ := json.Marshal(refreshReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		require.Equal(t, http.StatusOK, w.Code, "expected 200 OK on /refresh")

		var refreshRes apiResponse
		err := json.Unmarshal(w.Body.Bytes(), &refreshRes)
		require.NoError(t, err)

		rotatedAccessToken = refreshRes.Data["access_token"].(string)
		rotatedRefreshToken = refreshRes.Data["refresh_token"].(string)
		assert.NotEmpty(t, rotatedAccessToken)
		assert.NotEmpty(t, rotatedRefreshToken)
		assert.NotEqual(t, loginRefreshToken, rotatedRefreshToken, "refresh token must be rotated to a new token")
	})

	t.Run("5_ConcurrentRequest_GracePeriod", func(t *testing.T) {
		body, _ := json.Marshal(requests.RefreshTokenRequest{RefreshToken: loginRefreshToken})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		require.Equal(t, http.StatusOK, w.Code, "expected 200 OK for concurrent refresh within grace period")

		var res apiResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)

		concurrentAccessToken, ok := res.Data["access_token"].(string)
		require.True(t, ok && concurrentAccessToken != "")
		concurrentRefreshToken, ok := res.Data["refresh_token"].(string)
		require.True(t, ok && concurrentRefreshToken != "")
	})

	t.Run("6_TokenReuseDetection_Attack", func(t *testing.T) {
		// Age revoked_at past grace period to test replay attack detection without sleeping.
		loginHash := security.HashToken(loginRefreshToken)
		expiredGraceTime := time.Now().UTC().Add(-2 * time.Minute)
		err := db.Model(&models.RefreshToken{}).
			Where("token_hash = ?", loginHash).
			Update("revoked_at", expiredGraceTime).Error
		require.NoError(t, err)

		body, _ := json.Marshal(requests.RefreshTokenRequest{RefreshToken: loginRefreshToken})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "expected 401 Unauthorized on reusing old refresh token after grace period")

		rotatedHash := security.HashToken(rotatedRefreshToken)
		var rotatedRecord models.RefreshToken
		err = db.Where("token_hash = ?", rotatedHash).First(&rotatedRecord).Error
		require.NoError(t, err)
		assert.True(t, rotatedRecord.IsRevoked, "all user tokens must be revoked upon attack detection")
	})

	t.Run("7_Logout", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Identifier: testEmail,
			Password:   testPassword,
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := executeRequest(app, req)
		require.Equal(t, http.StatusOK, w.Code)

		var loginRes apiResponse
		_ = json.Unmarshal(w.Body.Bytes(), &loginRes)
		freshAccessToken := loginRes.Data["access_token"].(string)
		freshRefreshToken := loginRes.Data["refresh_token"].(string)

		logoutReq := requests.LogoutRequest{
			RefreshToken: freshRefreshToken,
		}
		logoutBody, _ := json.Marshal(logoutReq)
		req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewBuffer(logoutBody))
		req.Header.Set("Authorization", "Bearer "+freshAccessToken)
		req.Header.Set("Content-Type", "application/json")

		w = executeRequest(app, req)
		assert.Equal(t, http.StatusOK, w.Code, "expected 200 OK on logout")

		refreshBody, _ := json.Marshal(requests.RefreshTokenRequest{RefreshToken: freshRefreshToken})
		req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(refreshBody))
		req.Header.Set("Content-Type", "application/json")

		w = executeRequest(app, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "logged-out token must not have a grace period")
	})
}
