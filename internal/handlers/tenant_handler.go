package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"github.com/mubarok-ridho/casheer-auth-service/internal/utils"
)

type TenantHandler struct {
	DB *gorm.DB
}

func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{DB: db}
}

// GetProfile - mengambil profile tenant
func (h *TenantHandler) GetProfile(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	var tenant models.Tenant
	if err := h.DB.Preload("Users").First(&tenant, tenantID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	return c.JSON(tenant)
}

// SetupStore - mengatur toko (pertama kali atau update)
func (h *TenantHandler) SetupStore(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	var input struct {
		StoreName    string `json:"store_name"`
		StorePhone   string `json:"store_phone"`
		StoreEmail   string `json:"store_email"`
		StoreAddress string `json:"store_address"`
		ReceiptWidth string `json:"receipt_width"` // "58mm" atau "80mm"
		PrinterMAC   string `json:"printer_mac"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Mulai transaction
	tx := h.DB.Begin()

	// Update tenant
	var tenant models.Tenant
	if err := tx.First(&tenant, tenantID).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	tenant.StoreName = input.StoreName
	tenant.StorePhone = input.StorePhone
	tenant.StoreEmail = input.StoreEmail
	tenant.StoreAddress = input.StoreAddress
	tenant.ReceiptWidth = input.ReceiptWidth

	if err := tx.Save(&tenant).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update tenant",
		})
	}

	// Update atau create store settings
	var settings models.StoreSetting
	err := tx.Where("tenant_id = ?", tenantID).First(&settings).Error

	if err != nil {
		// Create new
		settings = models.StoreSetting{
			TenantID:     tenantID,
			PrinterMAC:   input.PrinterMAC,
			PrinterWidth: input.ReceiptWidth,
		}
		if err := tx.Create(&settings).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create store settings",
			})
		}
	} else {
		// Update existing
		settings.PrinterMAC = input.PrinterMAC
		settings.PrinterWidth = input.ReceiptWidth
		if err := tx.Save(&settings).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update store settings",
			})
		}
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Store setup completed successfully",
		"tenant":  tenant,
	})
}

// UploadLogo - upload logo toko ke Cloudinary
func (h *TenantHandler) UploadLogo(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	// Get file from form
	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No logo file uploaded",
		})
	}

	// Validasi file (max 2MB, hanya gambar)
	if file.Size > 2*1024*1024 {
		return c.Status(400).JSON(fiber.Map{
			"error": "File too large. Max size 2MB",
		})
	}

	// Upload ke Cloudinary
	logoURL, publicID, err := utils.UploadToCloudinary(file, "logos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to upload logo: " + err.Error(),
		})
	}

	// Update tenant dengan logo baru
	var tenant models.Tenant
	if err := h.DB.First(&tenant, tenantID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Tenant not found",
		})
	}

	tenant.LogoURL = logoURL
	if err := h.DB.Save(&tenant).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update tenant logo",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Logo uploaded successfully",
		"logo_url":  logoURL,
		"public_id": publicID,
	})
}
