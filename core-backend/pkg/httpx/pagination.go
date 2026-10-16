package httpx

import (
	"math"
	"strings"
)

const (
	DefaultPage               = 1
	DefaultPageSize           = 20
	MaxPageSize               = 100
	DefaultOrder              = "desc"
	PaginationOrderAscending  = "asc"
	PaginationOrderDescending = "desc"
)

type PaginationQuery struct {
	Page     int    `json:"page" query:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" query:"page_size" validate:"omitempty,min=1,max=100"`
	SortBy   string `json:"sort_by" query:"sort_by" validate:"omitempty,max=50"`
	Order    string `json:"order" query:"order" validate:"omitempty,oneof=asc desc ASC DESC"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	} else if p.PageSize > MaxPageSize {
		p.PageSize = MaxPageSize
	}

	order := strings.ToLower(strings.TrimSpace(p.Order))
	if order != PaginationOrderAscending && order != PaginationOrderDescending {
		p.Order = DefaultOrder
	} else {
		p.Order = order
	}

	p.SortBy = strings.TrimSpace(p.SortBy)
}

func (p PaginationQuery) Offset() int {
	page := p.Page
	if page < 1 {
		page = DefaultPage
	}
	pageSize := p.PageSize
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return (page - 1) * pageSize
}

func (p PaginationQuery) Limit() int {
	pageSize := p.PageSize
	if pageSize < 1 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationMeta(total int64, page, pageSize int) PaginationMeta {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	return PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

type PaginatedData[T any] struct {
	Items      []T            `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

func NewPaginatedData[T any](items []T, total int64, page, pageSize int) PaginatedData[T] {
	if items == nil {
		items = make([]T, 0)
	}
	return PaginatedData[T]{
		Items:      items,
		Pagination: NewPaginationMeta(total, page, pageSize),
	}
}

func NewPaginatedResponse[T any](items []T, total int64, page, pageSize int) Response {
	return NewResponse(NewPaginatedData(items, total, page, pageSize))
}
