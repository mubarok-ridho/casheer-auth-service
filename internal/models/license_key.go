package models

import (
	"time"
	"gorm.io/gorm"
)

type LicenseKey struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Key       string         `json:"key" gorm:"uniqueIndex;not null"`
	IsUsed    bool           `json:"is_used" gorm:"default:false"`
	UsedBy    *uint          `json:"used_by"` // tenant ID
	UsedAt    *time.Time     `json:"used_at"`
	Notes     string         `json:"notes"` // catatan: nama pembeli, dll
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
