package model

import (
	"github.com/example/core-backend/pkg/constants/response_code"
)

var (
	ErrNotFound      = response_code.NotFoundTenant
	ErrInvalidInput  = response_code.InvalidPayload
	ErrAlreadyExists = response_code.TenantAlreadyExists
)