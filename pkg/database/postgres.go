package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	// Baca dari environment variable
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "casheer_db")
	sslmode := getEnv("DB_SSLMODE", "disable")
	timezone := getEnv("DB_TIMEZONE", "Asia/Jakarta")

	// Connection string tanpa database name (untuk create database)
	dsnWithoutDB := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s TimeZone=%s",
		host, port, user, password, sslmode, timezone)

	// Coba konek ke postgres default database
	defaultDB, err := gorm.Open(postgres.Open(dsnWithoutDB+" dbname=postgres"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %v", err)
	}

	// Cek apakah database sudah ada, jika belum buat
	var count int64
	defaultDB.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbname).Scan(&count)

	if count == 0 {
		log.Printf("📦 Database %s not found, creating...", dbname)

		if err := defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname)).Error; err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}

		log.Printf("✅ Database %s created successfully", dbname)
	}

	// Tutup koneksi ke postgres database
	sqlDB, _ := defaultDB.DB()
	sqlDB.Close()

	// Tunggu sebentar agar database siap
	time.Sleep(2 * time.Second)

	// Konek ke database casheer_db
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %v", dbname, err)
	}

	// Set connection pool
	sqlDB, err = db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("✅ Connected to database: %s", dbname)
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
