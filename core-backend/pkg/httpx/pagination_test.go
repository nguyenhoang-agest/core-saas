package httpx_test

import (
	"testing"

	"github.com/example/core-backend/pkg/httpx"
)

func TestPaginationQuery_Normalize(t *testing.T) {
	tests := []struct {
		name         string
		input        httpx.PaginationQuery
		expectedPage int
		expectedSize int
		expectedOrd  string
	}{
		{
			name:         "Empty struct uses defaults",
			input:        httpx.PaginationQuery{},
			expectedPage: httpx.DefaultPage,
			expectedSize: httpx.DefaultPageSize,
			expectedOrd:  httpx.DefaultOrder,
		},
		{
			name: "Negative values are reset to defaults",
			input: httpx.PaginationQuery{
				Page:     -5,
				PageSize: -10,
				Order:    "invalid",
			},
			expectedPage: httpx.DefaultPage,
			expectedSize: httpx.DefaultPageSize,
			expectedOrd:  httpx.DefaultOrder,
		},
		{
			name: "PageSize exceeding max is capped",
			input: httpx.PaginationQuery{
				Page:     2,
				PageSize: 500,
				Order:    "asc",
			},
			expectedPage: 2,
			expectedSize: httpx.MaxPageSize,
			expectedOrd:  "asc",
		},
		{
			name: "Valid params are preserved",
			input: httpx.PaginationQuery{
				Page:     3,
				PageSize: 15,
				Order:    "ASC",
			},
			expectedPage: 3,
			expectedSize: 15,
			expectedOrd:  "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.input
			q.Normalize()

			if q.Page != tt.expectedPage {
				t.Errorf("expected page %d, got %d", tt.expectedPage, q.Page)
			}
			if q.PageSize != tt.expectedSize {
				t.Errorf("expected page_size %d, got %d", tt.expectedSize, q.PageSize)
			}
			if q.Order != tt.expectedOrd {
				t.Errorf("expected order %s, got %s", tt.expectedOrd, q.Order)
			}
		})
	}
}

func TestPaginationQuery_OffsetLimit(t *testing.T) {
	q := httpx.PaginationQuery{Page: 1, PageSize: 20}
	if q.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", q.Offset())
	}
	if q.Limit() != 20 {
		t.Errorf("expected limit 20, got %d", q.Limit())
	}

	q2 := httpx.PaginationQuery{Page: 3, PageSize: 15}
	if q2.Offset() != 30 {
		t.Errorf("expected offset 30, got %d", q2.Offset())
	}
	if q2.Limit() != 15 {
		t.Errorf("expected limit 15, got %d", q2.Limit())
	}
}

func TestNewPaginationMeta(t *testing.T) {
	tests := []struct {
		name          string
		total         int64
		page          int
		pageSize      int
		expectedPages int
	}{
		{"Zero total", 0, 1, 20, 0},
		{"Exact division", 100, 1, 20, 5},
		{"Ceil division", 101, 2, 20, 6},
		{"Single item", 1, 1, 20, 1},
		{"Less than one page", 15, 1, 20, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := httpx.NewPaginationMeta(tt.total, tt.page, tt.pageSize)
			if meta.TotalPages != tt.expectedPages {
				t.Errorf("expected total_pages %d, got %d", tt.expectedPages, meta.TotalPages)
			}
			if meta.TotalItems != tt.total {
				t.Errorf("expected total_items %d, got %d", tt.total, meta.TotalItems)
			}
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	items := []string{"item1", "item2"}
	resp := httpx.NewPaginatedResponse(items, 2, 1, 10)

	if resp.Status != "success" {
		t.Fatalf("expected status success, got %s", resp.Status)
	}

	data, ok := resp.Data.(httpx.PaginatedData[string])
	if !ok {
		t.Fatalf("expected PaginatedData[string], got %T", resp.Data)
	}

	if len(data.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(data.Items))
	}
	if data.Pagination.TotalItems != 2 {
		t.Errorf("expected total 2, got %d", data.Pagination.TotalItems)
	}
}
