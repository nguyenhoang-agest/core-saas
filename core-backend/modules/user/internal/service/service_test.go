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

func (m *mockRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	if item, ok := m.items[id]; ok {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) GetByName(ctx context.Context, name string) (*model.User, error) {
	for _, item := range m.items {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) Update(ctx context.Context, user *model.User) error {
	m.items[user.Id] = user
	return nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.items[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.items, id)
	return nil
}

func (m *mockRepo) GetAll(ctx context.Context, req *dto.PaginationOptions) ([]model.User, int64, error) {
	var list []model.User
	for _, item := range m.items {
		list = append(list, *item)
	}
	return list, int64(len(list)), nil
}

func TestUserService_Create(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	created, err := svc.Create(context.Background(), &dto.CreateUserRequest{
		Name:        "Test Item",
		Description: "A test description",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Name != "Test Item" {
		t.Fatalf("expected name 'Test Item', got '%s'", created.Name)
	}
	if created.Status != model.StatusActive {
		t.Fatalf("expected active status, got '%s'", created.Status)
	}
}

func TestUserService_Create_Duplicate(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	_, _ = svc.Create(context.Background(), &dto.CreateUserRequest{Name: "Existing Item"})
	_, err := svc.Create(context.Background(), &dto.CreateUserRequest{Name: "Existing Item"})

	var badReq exceptions.BadRequestError
	if !errors.As(err, &badReq) {
		t.Fatalf("expected BadRequestError, got %v", err)
	}
	if badReq.MessageCode != errcode.UserAlreadyExists.Code() {
		t.Fatalf("expected code %s, got %s", errcode.UserAlreadyExists.Code(), badReq.MessageCode)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := service.NewUserService(repo, nil)

	_, err := svc.GetByID(context.Background(), "non-existent")
	var notFound exceptions.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
	if notFound.MessageCode != errcode.NotFoundUser.Code() {
		t.Fatalf("expected code %s, got %s", errcode.NotFoundUser.Code(), notFound.MessageCode)
	}
}
