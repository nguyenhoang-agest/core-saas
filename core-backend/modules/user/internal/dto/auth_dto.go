package dto

import "github.com/example/core-backend/modules/user/internal/model"

type RegisterOrgRequest struct {
	OrgName  string `json:"org_name" validate:"required,min=3"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	TenantID string `json:"tenant_id,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthResponse struct {
	User         *model.User `json:"user"`
	TenantID     string      `json:"tenant_id"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
}

type TenantSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IsHost bool   `json:"is_host"`
}

type CurrentUserResponse struct {
	User        *model.User    `json:"user"`
	Tenant      *TenantSummary `json:"tenant"`
	Permissions []string       `json:"permissions"`
}
