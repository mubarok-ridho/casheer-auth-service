package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/mubarok-ridho/casheer-auth-service/internal/handlers"
	"github.com/mubarok-ridho/casheer-auth-service/internal/middleware"
	"github.com/mubarok-ridho/casheer-auth-service/internal/models"
	"github.com/mubarok-ridho/casheer-auth-service/internal/repository"
	"github.com/mubarok-ridho/casheer-auth-service/pkg/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found, using environment variables")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("📦 Running database migrations...")
	if err := db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.StoreSetting{},
		&models.LicenseKey{},
	); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	log.Println("✅ Database migration completed")

	app := fiber.New(fiber.Config{AppName: os.Getenv("APP_NAME")})
	app.Use(cors.New())
	setupRoutes(app, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)
	tenantHandler := handlers.NewTenantHandler(db)
	settingHandler := handlers.NewStoreSettingHandler(db)
	licenseRepo := repository.NewLicenseRepository(db)
	licenseHandler := handlers.NewLicenseHandler(licenseRepo)

	// Public
	app.Post("/api/v1/login", authHandler.Login)
	app.Post("/api/v1/register", authHandler.Register)

	// Admin internal
	admin := app.Group("/api/v1/admin", middleware.AdminMiddleware())
	admin.Post("/license/generate", licenseHandler.Generate)
	admin.Get("/license", licenseHandler.List)

	// License activate
	app.Post("/api/v1/license/activate", middleware.AuthMiddleware(), licenseHandler.Activate)

	// Protected
	api := app.Group("/api/v1", middleware.AuthMiddleware(), middleware.LicenseMiddleware(db))
	api.Get("/tenant/profile", tenantHandler.GetProfile)
	api.Put("/tenant/setup", tenantHandler.SetupStore)
	api.Post("/tenant/upload-logo", tenantHandler.UploadLogo)

	// Store settings
	api.Get("/settings", settingHandler.Get)
	api.Put("/settings", settingHandler.Update)
	api.Post("/settings/margins-password", settingHandler.SetMarginsPassword)
	api.Post("/settings/verify-margins-password", settingHandler.VerifyMarginsPassword)
}
