package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	StoreName       string         `json:"store_name" gorm:"not null"`
	StorePhone      string         `json:"store_phone"`
	StoreEmail      string         `json:"store_email"`
	StoreAddress    string         `json:"store_address"`
	LogoURL         string         `json:"logo_url" gorm:"not null"` // Wajib diisi
	ReceiptTemplate string         `json:"receipt_template" gorm:"default:'default'"`
	ReceiptWidth    string         `json:"receipt_width" gorm:"default:'58mm'"` // 58mm atau 80mm
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Users         []User         `json:"users,omitempty"`
	StoreSettings []StoreSetting `json:"store_settings,omitempty"`
}
