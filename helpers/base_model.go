package helpers

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	Id        int            `gorm:"column:id" json:"id"`
	CreatedAt int            `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
