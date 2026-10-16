package db

import "gorm.io/gorm"

// Paginate returns a GORM scope for offset/limit pagination based on page and pageSize.
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 20
		} else if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return tx.Offset(offset).Limit(pageSize)
	}
}

// PaginateOffset returns a GORM scope based on explicit offset and limit.
func PaginateOffset(offset, limit int) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if offset < 0 {
			offset = 0
		}
		if limit <= 0 {
			limit = 20
		} else if limit > 100 {
			limit = 100
		}
		return tx.Offset(offset).Limit(limit)
	}
}
