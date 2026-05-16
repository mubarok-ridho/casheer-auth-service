package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
)

type StoreSettingHandler struct {
	DB *gorm.DB
}

func NewStoreSettingHandler(db *gorm.DB) *StoreSettingHandler {
	return &StoreSettingHandler{DB: db}
}

func (h *StoreSettingHandler) getOrCreate(tenantID uint) (*models.StoreSetting, error) {
	var s models.StoreSetting
	err := h.DB.Where("tenant_id = ?", tenantID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		s = models.StoreSetting{
			TenantID:     tenantID,
			PrinterWidth: "58mm",
			Currency:     "IDR",
			LowStockAlert: 5,
		}
		if err := h.DB.Create(&s).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &s, nil
}

// GET /api/v1/settings
func (h *StoreSettingHandler) Get(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)
	s, err := h.getOrCreate(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

// PUT /api/v1/settings
func (h *StoreSettingHandler) Update(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	var input struct {
		PrinterMAC    string  `json:"printer_mac"`
		PrinterWidth  string  `json:"printer_width"`
		TaxRate       float64 `json:"tax_rate"`
		Currency      string  `json:"currency"`
		ReceiptHeader string  `json:"receipt_header"`
		ReceiptFooter string  `json:"receipt_footer"`
		EnhancedMode  bool    `json:"enhanced_mode"`
		LowStockAlert int     `json:"low_stock_alert"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	s, err := h.getOrCreate(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.PrinterMAC = input.PrinterMAC
	s.PrinterWidth = input.PrinterWidth
	s.TaxRate = input.TaxRate
	s.Currency = input.Currency
	s.ReceiptHeader = input.ReceiptHeader
	s.ReceiptFooter = input.ReceiptFooter
	s.EnhancedMode = input.EnhancedMode
	if input.LowStockAlert > 0 {
		s.LowStockAlert = input.LowStockAlert
	}

	if err := h.DB.Save(&s).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

// POST /api/v1/settings/margins-password
func (h *StoreSettingHandler) SetMarginsPassword(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	var input struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password wajib diisi"})
	}
	if len(input.Password) < 4 {
		return c.Status(400).JSON(fiber.Map{"error": "Password minimal 4 karakter"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
	}

	s, err := h.getOrCreate(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	s.MarginsPassword = string(hash)
	if err := h.DB.Save(&s).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Password margins berhasil diset"})
}

// POST /api/v1/settings/verify-margins-password
func (h *StoreSettingHandler) VerifyMarginsPassword(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(uint)

	var input struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password wajib diisi"})
	}

	s, err := h.getOrCreate(tenantID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if s.MarginsPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password margins belum diset"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(s.MarginsPassword), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Password salah"})
	}

	return c.JSON(fiber.Map{"verified": true})
}
