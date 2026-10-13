package dto

import "github.com/example/core-backend/modules/tenant/internal/model"

type CreateTenantRequest struct {
	Name        string                `json:"name" validate:"required,min=3,max=100"`
	Description string                `json:"description,omitempty" validate:"max=200"`
	IsolateMode model.IsolateMode     `json:"isolate_mode,omitempty" validate:"omitempty,oneof=shared schema dedicated"`
	DBConfig    *model.TenantDBConfig `json:"db_config,omitempty"`
}
