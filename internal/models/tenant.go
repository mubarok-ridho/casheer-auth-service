package models

import (
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	StoreName       string         `json:"store_name"`
	StorePhone      string         `json:"store_phone"`
	StoreEmail      string         `json:"store_email"`
	StoreAddress    string         `json:"store_address"`
	LogoURL         string         `json:"logo_url"` // Wajib diisi
	ReceiptTemplate string         `json:"receipt_template" gorm:"default:'default'"`
	ReceiptWidth    string         `json:"receipt_width" gorm:"default:'58mm'"` // 58mm atau 80mm
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Users           []User         `json:"users,omitempty"`
}
