package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	JWTIssuer        = "core-saas-iam"
)

var (
	ErrInvalidToken = errors.New("invalid or malformed token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenClaims represents standard and custom claims encoded in JWT tokens.
type TokenClaims struct {
	jwt.RegisteredClaims
	TenantID  string `json:"tenant_id"`
	Email     string `json:"email"`
	Role      string `json:"role,omitempty"`
	IsHost    bool   `json:"is_host"`
	TokenType string `json:"token_type"`
}

// TokenPair contains generated access and refresh tokens.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int64     `json:"expires_in"` // seconds
}

// GenerateTokenPair issues a short-lived access token and a long-lived refresh token.
func GenerateTokenPair(
	secret string,
	userID string,
	tenantID string,
	email string,
	role string,
	isHost bool,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) (*TokenPair, error) {
	if secret == "" {
		return nil, errors.New("jwt secret cannot be empty")
	}
	if accessTTL == 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL == 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	now := time.Now().UTC()
	accessExpiry := now.Add(accessTTL)
	refreshExpiry := now.Add(refreshTTL)

	// Access Token
	accessClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    JWTIssuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		IsHost:    isHost,
		TokenType: TokenTypeAccess,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh Token
	refreshClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    JWTIssuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		IsHost:    isHost,
		TokenType: TokenTypeRefresh,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  signedAccess,
		RefreshToken: signedRefresh,
		ExpiresAt:    accessExpiry,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

// ValidateToken parses and verifies a JWT token string with the given secret.
func ValidateToken(secret string, tokenStr string) (*TokenClaims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret cannot be empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
