package auth_test

import (
	"testing"
	"time"

	"github.com/example/core-backend/pkg/auth"
)

func TestPasswordHashing(t *testing.T) {
	plain := "P@ssword123!"

	hash, err := auth.HashPassword(plain)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == plain {
		t.Fatal("hashed password must not match plain text")
	}

	if !auth.CheckPassword(hash, plain) {
		t.Fatal("expected password check to succeed with correct password")
	}

	if auth.CheckPassword(hash, "WrongPassword!") {
		t.Fatal("expected password check to fail with wrong password")
	}
}

func TestJWTTokenGenerationAndValidation(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := "user-abc-123"
	tenantID := "tenant-xyz-789"
	email := "alice@acme.com"
	role := "owner"
	isHost := true

	tokens, err := auth.GenerateTokenPair(secret, userID, tenantID, email, role, isHost, 5*time.Minute, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("expected non-empty tokens in pair")
	}

	// Validate access token
	claims, err := auth.ValidateToken(secret, tokens.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("expected subject %s, got %s", userID, claims.Subject)
	}
	if claims.TenantID != tenantID {
		t.Errorf("expected tenantID %s, got %s", tenantID, claims.TenantID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
	if claims.IsHost != isHost {
		t.Errorf("expected isHost %v, got %v", isHost, claims.IsHost)
	}
	if claims.TokenType != auth.TokenTypeAccess {
		t.Errorf("expected tokenType %s, got %s", auth.TokenTypeAccess, claims.TokenType)
	}

	// Validate refresh token
	refreshClaims, err := auth.ValidateToken(secret, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}
	if refreshClaims.TokenType != auth.TokenTypeRefresh {
		t.Errorf("expected tokenType %s, got %s", auth.TokenTypeRefresh, refreshClaims.TokenType)
	}

	// Validation with wrong secret
	_, err = auth.ValidateToken("wrong-secret", tokens.AccessToken)
	if err == nil {
		t.Fatal("expected token validation to fail with wrong secret")
	}
}

func TestJWTTokenExpiration(t *testing.T) {
	secret := "test-secret-key-12345"
	// Expired token (negative TTL)
	tokens, err := auth.GenerateTokenPair(secret, "uid", "tid", "a@b.com", "member", false, -1*time.Minute, -1*time.Minute)
	if err != nil {
		t.Fatalf("failed to generate tokens: %v", err)
	}

	_, err = auth.ValidateToken(secret, tokens.AccessToken)
	if err != auth.ErrExpiredToken {
		t.Fatalf("expected ErrExpiredToken, got %v", err)
	}
}
