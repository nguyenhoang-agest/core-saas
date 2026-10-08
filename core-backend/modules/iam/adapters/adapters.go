package adapters

import (
	"context"

	"github.com/example/core-backend/modules/iam/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresUserRepository implements domain.Repository (STUB).
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	// STUB: intentionally omitted
	return errors.Internal("not implemented", nil)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	// STUB: intentionally omitted
	return nil, errors.NotFound("user not found", nil)
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	// STUB: intentionally omitted
	return nil, errors.NotFound("user not found", nil)
}

// BcryptHasher implements domain.PasswordHasher (STUB).
type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	// STUB: replace with bcrypt.GenerateFromPassword in production
	return "stub-hash:" + password, nil
}

func (h *BcryptHasher) Compare(hashed, plain string) bool {
	// STUB: replace with bcrypt.CompareHashAndPassword in production
	return hashed == "stub-hash:"+plain
}

// TokenStub implements domain.TokenGenerator (STUB).
type TokenStub struct{}

func NewTokenStub() *TokenStub {
	return &TokenStub{}
}

func (t *TokenStub) Generate(user *domain.User) (string, error) {
	return "stub-jwt-token-for-" + user.ID, nil
}

func (t *TokenStub) Validate(token string) (string, error) {
	if token == "" {
		return "", errors.Unauthorized("missing token", nil)
	}
	return "stub-user-id", nil
}
