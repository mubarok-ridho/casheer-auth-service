package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TenantID  uint           `json:"tenant_id" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-"`                             // "-" agar tidak muncul di JSON
	Role      string         `json:"role" gorm:"default:'cashier'"` // admin, cashier
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}
