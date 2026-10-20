package dto

import (
	"time"

	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/httpx"
)

type CreateUserRequest struct {
	Email       string       `json:"email" validate:"required,email,max=255"`
	Name        string       `json:"name" validate:"required,min=2,max=100"`
	Role        string       `json:"role" validate:"required,oneof=owner admin member viewer"`
	Permissions []string     `json:"permissions,omitempty"`
	Status      model.Status `json:"status,omitempty" validate:"omitempty,oneof=active suspended invited"`
	Phone       string       `json:"phone,omitempty" validate:"omitempty,max=50"`
}

type UpdateUserRequest struct {
	Name        string       `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Role        string       `json:"role,omitempty" validate:"omitempty,oneof=owner admin member viewer"`
	Permissions *[]string    `json:"permissions,omitempty"`
	Status      model.Status `json:"status,omitempty" validate:"omitempty,oneof=active suspended invited"`
	Phone       string       `json:"phone,omitempty" validate:"omitempty,max=50"`
}

type UpdateUserStatusRequest struct {
	Status model.Status `json:"status" validate:"required,oneof=active suspended invited"`
}

type PaginationOptions struct {
	httpx.PaginationQuery
	Keyword string       `json:"keyword" query:"keyword" validate:"omitempty,max=100"`
	Role    string       `json:"role" query:"role" validate:"omitempty,oneof=owner admin member viewer"`
	Status  model.Status `json:"status" query:"status" validate:"omitempty,oneof=active suspended invited"`
}

type UserResponse struct {
	Id                   string       `json:"id"`
	TenantID             string       `json:"tenant_id"`
	Email                string       `json:"email"`
	Name                 string       `json:"name"`
	Role                 string       `json:"role"`
	Permissions          []string     `json:"permissions"`
	EffectivePermissions []string     `json:"effective_permissions"`
	Status               model.Status `json:"status"`
	Phone                string       `json:"phone,omitempty"`
	AvatarURL            string       `json:"avatar_url,omitempty"`
	LastLoginAt          *time.Time   `json:"last_login_at,omitempty"`
	CreateAt             time.Time    `json:"create_at"`
	UpdateAt             time.Time    `json:"update_at"`
}

func ToUserResponse(u *model.User) *UserResponse {
	if u == nil {
		return nil
	}
	perms := []string(u.Permissions)
	if perms == nil {
		perms = []string{}
	}
	return &UserResponse{
		Id:                   u.Id,
		TenantID:             u.TenantID,
		Email:                u.Email,
		Name:                 u.Name,
		Role:                 u.Role,
		Permissions:          perms,
		EffectivePermissions: u.EffectivePermissions(),
		Status:               u.Status,
		Phone:                u.Phone,
		AvatarURL:            u.AvatarURL,
		LastLoginAt:          u.LastLoginAt,
		CreateAt:             u.CreateAt,
		UpdateAt:             u.UpdateAt,
	}
}
