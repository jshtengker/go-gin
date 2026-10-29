package security

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAccessTokenGenerationAndValidation(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	username := "testuser"
	secret := "test-secret-key-12345"
	expiry := 15 * time.Minute

	token, err := GenerateAccessToken(userID, email, username, secret, expiry)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	claims, err := ValidateAccessToken(token, secret)
	if err != nil {
		t.Fatalf("failed to validate valid access token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}
	if claims.Username != username {
		t.Errorf("expected Username %s, got %s", username, claims.Username)
	}

	_, err = ValidateAccessToken(token, "wrong-secret")
	if !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("expected ErrTokenInvalid when validating with wrong secret, got: %v", err)
	}
}

func TestAccessTokenExpired(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key-12345"

	token, err := GenerateAccessToken(userID, "test@example.com", "testuser", secret, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate expired access token: %v", err)
	}

	_, err = ValidateAccessToken(token, secret)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired for expired token, got: %v", err)
	}
}

func TestRefreshTokenGenerationAndHashing(t *testing.T) {
	rawToken, tokenHash, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}

	if len(rawToken) == 0 {
		t.Fatal("expected non-empty raw refresh token")
	}

	if len(tokenHash) != 64 {
		t.Fatalf("expected 64-char sha256 hex hash, got length %d", len(tokenHash))
	}

	computedHash := HashToken(rawToken)
	if computedHash != tokenHash {
		t.Fatalf("expected HashToken to produce %s, got %s", tokenHash, computedHash)
	}
}
