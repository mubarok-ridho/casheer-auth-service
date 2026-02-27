package handlers

import (
	"casheer-auth-service/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	if err := h.DB.Preload("Tenant").Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
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

// Setup Toko handler
func (h *AuthHandler) SetupStore(c *fiber.Ctx) error {
	// Dapatkan tenant_id dari JWT (via middleware)
	tenantID := c.Locals("tenant_id").(uint)

	var input struct {
		StoreName    string `json:"store_name"`
		StorePhone   string `json:"store_phone"`
		StoreEmail   string `json:"store_email"`
		StoreAddress string `json:"store_address"`
		ReceiptWidth string `json:"receipt_width"` // "58mm" atau "80mm"
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var tenant models.Tenant
	if err := h.DB.First(&tenant, tenantID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tenant not found"})
	}

	// Update tenant
	tenant.StoreName = input.StoreName
	tenant.StorePhone = input.StorePhone
	tenant.StoreEmail = input.StoreEmail
	tenant.StoreAddress = input.StoreAddress
	tenant.ReceiptWidth = input.ReceiptWidth

	// Handle logo upload (dari multipart form)
	file, err := c.FormFile("logo")
	if err == nil {
		// Upload ke Cloudinary
		// Implementasi upload ke Cloudinary
		// tenant.LogoURL = uploadedURL
	}

	h.DB.Save(&tenant)

	return c.JSON(tenant)
}
