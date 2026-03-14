package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Login handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		LicenseKey string `json:"license_key"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Validasi input
	if input.Email == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	var user models.User
	err := h.DB.Preload("Tenant").Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(401).JSON(fiber.Map{
			"error": "Account is deactivated",
		})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"tenant_id":  user.TenantID,
			"store_name": user.Tenant.StoreName,
			"logo_url":   user.Tenant.LogoURL,
		},
	})
}

// Register handler (untuk registrasi tenant baru)
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input struct {
		StoreName  string `json:"store_name"`
		StorePhone string `json:"store_phone"`
		StoreEmail string `json:"store_email"`
		AdminName  string `json:"admin_name"`
		AdminEmail string `json:"admin_email"`
		Password   string `json:"password"`
		LicenseKey string `json:"license_key"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Validasi
	if input.StoreName == "" || input.AdminEmail == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Store name, admin email and password are required",
		})
	}
	// Validasi license key
	if input.LicenseKey == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "License key diperlukan",
		})
	}
	var license models.LicenseKey
	if err := h.DB.Where("key = ? AND is_used = false", strings.ToUpper(strings.TrimSpace(input.LicenseKey))).First(&license).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "License key tidak valid atau sudah digunakan",
		})
	}

	// Start transaction
	tx := h.DB.Begin()

	// Buat tenant baru
	tenant := models.Tenant{
		StoreName:  input.StoreName,
		StorePhone: input.StorePhone,
		StoreEmail: input.StoreEmail,
		LogoURL:    "https://res.cloudinary.com/demo/image/upload/v1/default/logo.png", // Default logo
	}

	if err := tx.Create(&tenant).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create tenant: " + err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Buat admin user
	user := models.User{
		TenantID: tenant.ID,
		Name:     input.AdminName,
		Email:    input.AdminEmail,
		Password: string(hashedPassword),
		Role:     "admin",
		IsActive: true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user: " + err.Error(),
		})
	}

	// Buat default store settings
	settings := models.StoreSetting{
		TenantID:     tenant.ID,
		PrinterWidth: "58mm",
		TaxRate:      0,
		Currency:     "IDR",
	}

	if err := tx.Create(&settings).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create store settings",
		})
	}

	tx.Commit()

	// Aktivasi license
	now := time.Now()
	h.DB.Model(&license).Updates(map[string]interface{}{
		"is_used": true,
		"used_by": tenant.ID,
		"used_at": now,
	})
	h.DB.Model(&tenant).Updates(map[string]interface{}{
		"is_active":   true,
		"license_key": strings.ToUpper(strings.TrimSpace(input.LicenseKey)),
	})

	return c.Status(201).JSON(fiber.Map{
		"message":   "Tenant registered successfully",
		"tenant_id": tenant.ID,
	})
}
