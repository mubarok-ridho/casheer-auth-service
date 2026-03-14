package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mubarok-ridho/casheer-auth-service/internal/repository"
)

type LicenseHandler struct {
	repo *repository.LicenseRepository
}

func NewLicenseHandler(repo *repository.LicenseRepository) *LicenseHandler {
	return &LicenseHandler{repo: repo}
}

func (h *LicenseHandler) Activate(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	var input struct {
		Key string `json:"key"`
	}
	if err := c.BodyParser(&input); err != nil || input.Key == "" {
		return c.Status(400).JSON(fiber.Map{"error": "License key diperlukan"})
	}
	if err := h.repo.Activate(tenantID, input.Key); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Aktivasi berhasil! Selamat menggunakan MODU."})
}

func (h *LicenseHandler) Generate(c *fiber.Ctx) error {
	var input struct {
		Count int    `json:"count"`
		Notes string `json:"notes"`
	}
	c.BodyParser(&input)
	if input.Count <= 0 {
		input.Count = 1
	}
	keys, err := h.repo.GenerateKeys(input.Count, input.Notes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"keys": keys, "count": len(keys)})
}

func (h *LicenseHandler) List(c *fiber.Ctx) error {
	keys, err := h.repo.ListKeys()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(keys)
}
