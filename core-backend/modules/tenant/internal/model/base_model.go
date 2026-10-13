package model

import (
	"time"
)

type BaseModel struct {
	Id         string `gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" json:"id"`
	AuditModel `gorm:"embedded"`
}

type AuditModel struct {
	CreateAt time.Time  `gorm:"column:create_at;autoCreateTime" json:"create_at"`
	UpdateAt time.Time  `gorm:"column:update_at;autoUpdateTime" json:"update_at"`
	DeleteAt *time.Time `gorm:"column:delete_at" json:"delete_at,omitempty"`
}
