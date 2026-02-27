package repository

import (
	"errors"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create user
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// Get by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Tenant").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Get by Email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Tenant").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Get users by tenant
func (r *UserRepository) GetByTenant(tenantID uint) ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("tenant_id = ?", tenantID).Find(&users).Error
	return users, err
}

// Update user
func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

// Delete user (soft delete)
func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&models.User{}, id).Error
}

// Update password
func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// Toggle active status
func (r *UserRepository) ToggleActive(userID uint) error {
	return r.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_active", gorm.Expr("NOT is_active")).Error
}
