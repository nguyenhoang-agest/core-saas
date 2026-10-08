package exception

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("tenant not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("tenant already exists")
)