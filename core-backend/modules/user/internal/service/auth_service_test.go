package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/exceptions"
)

func TestAuthService_RegisterOrg_Success(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	req := &dto.RegisterOrgRequest{
		OrgName:  "Acme Corporation",
		Name:     "Alice Owner",
		Email:    "alice@acme.com",
		Password: "Password123!",
	}

	res, err := authSvc.RegisterOrg(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error registering org: %v", err)
	}

	if res.TenantID == "" {
		t.Fatal("expected non-empty TenantID")
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatal("expected generated tokens")
	}
	if res.User.Email != "alice@acme.com" {
		t.Errorf("expected email alice@acme.com, got %s", res.User.Email)
	}
	if res.User.Role != "owner" {
		t.Errorf("expected role owner, got %s", res.User.Role)
	}

	// Verify user saved with hashed password
	savedUser, err := repo.GetByEmail(context.Background(), res.TenantID, "alice@acme.com")
	if err != nil {
		t.Fatalf("failed to retrieve saved user: %v", err)
	}
	if !auth.CheckPassword(savedUser.PasswordHash, "Password123!") {
		t.Fatal("saved user password hash did not verify")
	}
}

func TestAuthService_RegisterOrg_Validations(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	// Short org name
	_, err := authSvc.RegisterOrg(context.Background(), &dto.RegisterOrgRequest{
		OrgName:  "ab",
		Name:     "Alice",
		Email:    "a@b.com",
		Password: "Password123!",
	})
	if err == nil {
		t.Fatal("expected error for short org name")
	}

	// Short password
	_, err = authSvc.RegisterOrg(context.Background(), &dto.RegisterOrgRequest{
		OrgName:  "Acme Corp",
		Name:     "Alice",
		Email:    "a@b.com",
		Password: "short",
	})
	if err == nil {
		t.Fatal("expected error for short password")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	passwordHash, _ := auth.HashPassword("SecretPassword!")
	user := &model.User{
		BaseModel:    model.BaseModel{Id: "user-123"},
		TenantID:     "tenant-1",
		Email:        "bob@acme.com",
		Name:         "Bob",
		Role:         "admin",
		Status:       model.StatusActive,
		PasswordHash: passwordHash,
	}
	_ = repo.Create(context.Background(), user)

	res, err := authSvc.Login(context.Background(), &dto.LoginRequest{
		Email:    "bob@acme.com",
		Password: "SecretPassword!",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}

	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatal("expected token pair on login")
	}
	if res.User.Id != "user-123" {
		t.Errorf("expected user id user-123, got %s", res.User.Id)
	}

	// Verify claims
	claims, err := auth.ValidateToken(jwtSecret, res.AccessToken)
	if err != nil {
		t.Fatalf("invalid token returned: %v", err)
	}
	if claims.Subject != "user-123" || claims.TenantID != "tenant-1" {
		t.Errorf("claims mismatch: subject=%s tenant=%s", claims.Subject, claims.TenantID)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	passwordHash, _ := auth.HashPassword("CorrectPassword!")
	user := &model.User{
		BaseModel:    model.BaseModel{Id: "user-123"},
		TenantID:     "tenant-1",
		Email:        "bob@acme.com",
		Status:       model.StatusActive,
		PasswordHash: passwordHash,
	}
	_ = repo.Create(context.Background(), user)

	// Wrong password
	_, err := authSvc.Login(context.Background(), &dto.LoginRequest{
		Email:    "bob@acme.com",
		Password: "WrongPassword!",
	})
	if err == nil {
		t.Fatal("expected unauthorized error for wrong password")
	}

	// Non-existent email
	_, err = authSvc.Login(context.Background(), &dto.LoginRequest{
		Email:    "unknown@acme.com",
		Password: "SomePassword!",
	})
	if err == nil {
		t.Fatal("expected unauthorized error for unknown email")
	}
}

func TestAuthService_Login_SuspendedUser(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	passwordHash, _ := auth.HashPassword("CorrectPassword!")
	user := &model.User{
		BaseModel:    model.BaseModel{Id: "user-suspended"},
		TenantID:     "tenant-1",
		Email:        "suspended@acme.com",
		Status:       model.StatusSuspended,
		PasswordHash: passwordHash,
	}
	_ = repo.Create(context.Background(), user)

	_, err := authSvc.Login(context.Background(), &dto.LoginRequest{
		Email:    "suspended@acme.com",
		Password: "CorrectPassword!",
	})
	if err == nil {
		t.Fatal("expected forbidden error for suspended user")
	}
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	user := &model.User{
		BaseModel: model.BaseModel{Id: "user-refresh"},
		TenantID:  "tenant-1",
		Email:     "ref@acme.com",
		Status:    model.StatusActive,
	}
	_ = repo.Create(context.Background(), user)

	tokens, err := auth.GenerateTokenPair(jwtSecret, user.Id, user.TenantID, user.Email, user.Role, false, 15*time.Minute, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	refreshed, err := authSvc.RefreshToken(context.Background(), &dto.RefreshTokenRequest{
		RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		t.Fatalf("failed to refresh token: %v", err)
	}

	if refreshed.AccessToken == "" {
		t.Fatal("expected new access token")
	}
}

func TestAuthService_GetCurrentUser(t *testing.T) {
	repo := newMockRepo()
	jwtSecret := "test-jwt-secret-key-12345"
	authSvc := service.NewAuthService(nil, repo, nil, nil, jwtSecret, 15*time.Minute, 7*24*time.Hour)

	user := &model.User{
		BaseModel: model.BaseModel{Id: "user-me"},
		TenantID:  "tenant-1",
		Email:     "me@acme.com",
		Name:      "Me User",
		Role:      "admin",
		Status:    model.StatusActive,
	}
	_ = repo.Create(context.Background(), user)

	me, err := authSvc.GetCurrentUser(context.Background(), "tenant-1", "user-me", false)
	if err != nil {
		t.Fatalf("unexpected error in GetCurrentUser: %v", err)
	}

	if me.User.Id != "user-me" {
		t.Errorf("expected user-me, got %s", me.User.Id)
	}
	if me.Tenant.ID != "tenant-1" {
		t.Errorf("expected tenant-1, got %s", me.Tenant.ID)
	}

	// Unknown user should return 404
	_, err = authSvc.GetCurrentUser(context.Background(), "tenant-1", "non-existent", false)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
	var notFound *exceptions.NotFoundError
	if !errors.As(err, &notFound) && !isNotFoundError(err) {
		t.Errorf("expected NotFound error, got %v", err)
	}
}

func isNotFoundError(err error) bool {
	var target *exceptions.NotFoundError
	return errors.As(err, &target) || (err != nil && err.Error() == errcode.NotFoundUser.Code())
}
