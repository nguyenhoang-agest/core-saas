package dto

import (
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/pkg/httpx"
)

type CreateTenantRequest struct {
	Name        string                `json:"name" validate:"required,min=3,max=100"`
	Description string                `json:"description,omitempty" validate:"max=200"`
	IsolateMode model.IsolateMode     `json:"isolate_mode,omitempty" validate:"omitempty,oneof=shared schema dedicated"`
	DBConfig    *model.TenantDBConfig `json:"db_config,omitempty"`
}

type PaginationOptions struct {
	httpx.PaginationQuery
	Keyword     string            `json:"keyword" query:"keyword" validate:"omitempty,max=100"`
	Status      model.Status      `json:"status" query:"status" validate:"omitempty,oneof=active suspended"`
	IsolateMode model.IsolateMode `json:"isolate_mode" query:"isolate_mode" validate:"omitempty,oneof=shared dedicated schema"`
}
