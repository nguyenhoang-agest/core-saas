package domain

import (
	"context"
	"time"
)

// User represents an IAM identity.
type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Repository defines data access ports for User.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, tenantID, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

// PasswordHasher defines the port for hashing and verifying passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) bool
}

// TokenGenerator defines the port for generating and validating authentication tokens.
type TokenGenerator interface {
	Generate(user *User) (string, error)
	Validate(token string) (string, error)
}
