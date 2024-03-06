package models

import (
	"time"

	"gorm.io/gorm"
)

type TimeCommon struct {
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggerignore:"true"`
}
