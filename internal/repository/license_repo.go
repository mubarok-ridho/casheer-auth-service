package repository

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"gorm.io/gorm"
)

type LicenseRepository struct {
	DB *gorm.DB
}

func NewLicenseRepository(db *gorm.DB) *LicenseRepository {
	return &LicenseRepository{DB: db}
}

func generateKey() string {
	b := make([]byte, 12)
	rand.Read(b)
	raw := fmt.Sprintf("%X", b)
	return fmt.Sprintf("MODU-%s-%s-%s", raw[0:4], raw[4:8], raw[8:12])
}

func (r *LicenseRepository) GenerateKeys(count int, notes string) ([]models.LicenseKey, error) {
	var keys []models.LicenseKey
	for i := 0; i < count; i++ {
		key := models.LicenseKey{
			Key:   generateKey(),
			Notes: notes,
		}
		if err := r.DB.Create(&key).Error; err != nil {
			continue
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (r *LicenseRepository) Activate(tenantID uint, keyStr string) error {
	keyStr = strings.ToUpper(strings.TrimSpace(keyStr))
	var license models.LicenseKey
	if err := r.DB.Where("key = ?", keyStr).First(&license).Error; err != nil {
		return fmt.Errorf("license key tidak ditemukan")
	}
	if license.IsUsed {
		return fmt.Errorf("license key sudah digunakan")
	}
	now := time.Now()
	license.IsUsed = true
	license.UsedBy = &tenantID
	license.UsedAt = &now
	r.DB.Save(&license)
	return r.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Updates(map[string]interface{}{
		"is_active":   true,
		"license_key": keyStr,
	}).Error
}

func (r *LicenseRepository) ListKeys() ([]models.LicenseKey, error) {
	var keys []models.LicenseKey
	err := r.DB.Order("created_at desc").Find(&keys).Error
	return keys, err
}
