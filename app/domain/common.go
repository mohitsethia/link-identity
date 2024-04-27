package domain

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
