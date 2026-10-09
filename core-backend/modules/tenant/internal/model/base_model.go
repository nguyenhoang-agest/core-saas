package model

import (
	"time"
)

type BaseModel struct {
	Id string `gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" json:"id"`
	AuditModel
}

type AuditModel struct {
	CreateAt time.Time  `gorm:"column:create_at;type:datetime,autoCreateTime" json:"create_at"`
	UpdateAt time.Time  `gorm:"column:update_at;type:datetime,autoUpdateTime" json:"update_at"`
	DeleteAt *time.Time `gorm:"column:delete_at;type:datetime,autoDeleteTime" json:"delete_at,omitempty"`
}
