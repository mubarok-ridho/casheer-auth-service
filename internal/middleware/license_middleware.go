package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"gorm.io/gorm"
)

// LicenseMiddleware - cek apakah tenant sudah aktif
func LicenseMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Next()
		}

		var tenant models.Tenant
		if err := db.Select("is_active").First(&tenant, tenantID).Error; err != nil {
			return c.Next()
		}

		if !tenant.IsActive {
			return c.Status(403).JSON(fiber.Map{
				"error":   "license_required",
				"message": "Akun belum diaktifkan. Masukkan license key untuk melanjutkan.",
			})
		}

		return c.Next()
	}
}
