package repository

import (
	"errors"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"gorm.io/gorm"
)

type TenantRepository struct {
	DB *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{DB: db}
}

// Create tenant
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.DB.Create(tenant).Error
}

// Get by ID
func (r *TenantRepository) GetByID(id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.DB.Preload("Users").First(&tenant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// Update tenant
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.DB.Save(tenant).Error
}

// Delete tenant (soft delete)
func (r *TenantRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Tenant{}, id).Error
}

// Get all tenants (admin only)
func (r *TenantRepository) GetAll() ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.DB.Find(&tenants).Error
	return tenants, err
}

// Get tenant by store name (search)
func (r *TenantRepository) SearchByName(name string) ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := r.DB.Where("store_name ILIKE ?", "%"+name+"%").Find(&tenants).Error
	return tenants, err
}

// Update logo
func (r *TenantRepository) UpdateLogo(tenantID uint, logoURL string) error {
	return r.DB.Model(&models.Tenant{}).Where("id = ?", tenantID).Update("logo_url", logoURL).Error
}
