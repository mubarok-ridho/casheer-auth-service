package middleware

import (
	"os"
	"github.com/gofiber/fiber/v2"
)

// AdminMiddleware - cek header X-Admin-Key
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminKey := c.Get("X-Admin-Key")
		expectedKey := os.Getenv("ADMIN_SECRET_KEY")
		if expectedKey == "" {
			expectedKey = "modu-admin-secret"
		}
		if adminKey != expectedKey {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Next()
	}
}
