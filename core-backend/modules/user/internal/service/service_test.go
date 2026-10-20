package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/tenantctx"
	"gorm.io/gorm"
)

type mockRepo struct {
	items map[string]*model.User
}

func newMockRepo() *mockRepo {
	return &mockRepo{items: make(map[string]*model.User)}
}

func (m *mockRepo) Create(ctx context.Context, user *model.User) error {
	if user.Id == "" {
		user.Id = "id-1"
	}
	m.items[user.Id] = user
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, tenantID, id string) (*model.User, error) {
	if item, ok := m.items[id]; ok && item.TenantID == tenantID {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) GetByEmail(ctx context.Context, tenantID, email string) (*model.User, error) {
	for _, item := range m.items {
		if item.TenantID == tenantID && item.Email == email {
			return item, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) Update(ctx context.Context, user *model.User) error {
	m.items[user.Id] = user
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, tenantID, id string) error {
	if item, ok := m.items[id]; !ok || item.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *mockRepo) GetAll(ctx context.Context, tenantID string, req *dto.PaginationOptions) ([]model.User, int64, error) {
	var list []model.User
	for _, item := range m.items {
		if item.TenantID == tenantID {
			list = append(list, *item)
		}
	}
	return list, int64(len(list)), nil
}

func TestUserService_Create_Success(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	ctx := tenantctx.WithContext(context.Background(), "tenant-alpha")
	created, err := svc.Create(ctx, &dto.CreateUserRequest{
		Email:       "alice@example.com",
		Name:        "Alice Member",
		Role:        "member",
		Permissions: []string{"billing:view"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %s", created.Email)
	}
	if created.TenantID != "tenant-alpha" {
		t.Fatalf("expected tenant-alpha, got %s", created.TenantID)
	}
	if created.Role != "member" {
		t.Fatalf("expected role member, got %s", created.Role)
	}
	if len(created.Permissions) != 1 || created.Permissions[0] != "billing:view" {
		t.Fatalf("expected custom ACL permission billing:view, got %v", created.Permissions)
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	ctx := tenantctx.WithContext(context.Background(), "tenant-alpha")
	_, _ = svc.Create(ctx, &dto.CreateUserRequest{
		Email: "existing@example.com",
		Name:  "First User",
		Role:  "admin",
	})

	_, err := svc.Create(ctx, &dto.CreateUserRequest{
		Email: "existing@example.com",
		Name:  "Second User",
		Role:  "member",
	})

	var badReq exceptions.BadRequestError
	if !errors.As(err, &badReq) {
		t.Fatalf("expected BadRequestError, got %v", err)
	}
	if badReq.MessageCode != errcode.UserAlreadyExists.Code() {
		t.Fatalf("expected code %s, got %s", errcode.UserAlreadyExists.Code(), badReq.MessageCode)
	}
}

func TestUserService_Create_InvalidRole(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	ctx := tenantctx.WithContext(context.Background(), "tenant-alpha")
	_, err := svc.Create(ctx, &dto.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Bad Role User",
		Role:  "superman",
	})

	var badReq exceptions.BadRequestError
	if !errors.As(err, &badReq) {
		t.Fatalf("expected BadRequestError, got %v", err)
	}
	if badReq.MessageCode != errcode.UserInvalidRole.Code() {
		t.Fatalf("expected code %s, got %s", errcode.UserInvalidRole.Code(), badReq.MessageCode)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	ctx := tenantctx.WithContext(context.Background(), "tenant-alpha")
	_, err := svc.GetByID(ctx, "non-existent")
	var notFound exceptions.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
	if notFound.MessageCode != errcode.NotFoundUser.Code() {
		t.Fatalf("expected code %s, got %s", errcode.NotFoundUser.Code(), notFound.MessageCode)
	}
}

func TestUserService_Update_RoleAndACL(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	ctx := tenantctx.WithContext(context.Background(), "tenant-alpha")
	created, _ := svc.Create(ctx, &dto.CreateUserRequest{
		Email: "bob@example.com",
		Name:  "Bob Member",
		Role:  "member",
	})

	newPerms := []string{"settings:branding:update"}
	updated, err := svc.Update(ctx, created.Id, &dto.UpdateUserRequest{
		Role:        "admin",
		Permissions: &newPerms,
	})
	if err != nil {
		t.Fatalf("unexpected error updating user: %v", err)
	}

	if updated.Role != "admin" {
		t.Fatalf("expected role admin, got %s", updated.Role)
	}
	if len(updated.Permissions) != 1 || updated.Permissions[0] != "settings:branding:update" {
		t.Fatalf("expected custom ACL permission, got %v", updated.Permissions)
	}
}
