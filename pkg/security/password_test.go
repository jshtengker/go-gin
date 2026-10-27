package security

import (
	"testing"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "mySecretPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty password hash")
	}

	if !CheckPassword(password, hash) {
		t.Fatal("expected CheckPassword to return true for matching password")
	}

	if CheckPassword("wrongPassword", hash) {
		t.Fatal("expected CheckPassword to return false for mismatched password")
	}
}
