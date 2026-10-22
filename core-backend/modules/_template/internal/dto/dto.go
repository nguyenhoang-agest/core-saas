package dto

import (
	"github.com/example/core-backend/modules/_template/internal/model"
	"github.com/example/core-backend/pkg/httpx"
)

type CreateTemplateRequest struct {
	Name        string       `json:"name" validate:"required,min=3,max=100"`
	Description string       `json:"description,omitempty" validate:"max=200"`
	Status      model.Status `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

type UpdateTemplateRequest struct {
	Name        string       `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description string       `json:"description,omitempty" validate:"omitempty,max=200"`
	Status      model.Status `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

type PaginationOptions struct {
	httpx.PaginationQuery
	Keyword string       `json:"keyword" query:"keyword" validate:"omitempty,max=100"`
	Status  model.Status `json:"status" query:"status" validate:"omitempty,oneof=active inactive"`
}
