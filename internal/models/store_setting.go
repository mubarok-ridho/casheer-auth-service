package models

import (
	"time"

	"gorm.io/gorm"
)

type StoreSetting struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	TenantID      uint           `json:"tenant_id" gorm:"not null;uniqueIndex"`
	PrinterMAC    string         `json:"printer_mac"` // Alamat MAC printer bluetooth
	PrinterWidth  string         `json:"printer_width" gorm:"default:'58mm'"`
	TaxRate       float64        `json:"tax_rate" gorm:"default:0"`
	Currency      string         `json:"currency" gorm:"default:'IDR'"`
	ReceiptHeader string         `json:"receipt_header"`
	ReceiptFooter string         `json:"receipt_footer"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}
