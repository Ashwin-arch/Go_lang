package database

import (
	"fmt"
	"log"

	"banking-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB connects to SQLite database and runs auto-migrations
func InitDB(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🏦 Banking SQLite Database connected.")

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&models.Account{},
		&models.Transaction{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	SeedData(DB)
	return DB
}

// SeedData initializes initial demo bank accounts
func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&models.Account{}).Count(&count)
	if count > 0 {
		return
	}

	db.Create(&models.Account{
		AccountNumber: "ACC1001",
		AccountHolder: "Alice Smith",
		Balance:       5000.00,
		Currency:      "USD",
		Status:        "ACTIVE",
	})

	db.Create(&models.Account{
		AccountNumber: "ACC1002",
		AccountHolder: "Bob Johnson",
		Balance:       1200.00,
		Currency:      "USD",
		Status:        "ACTIVE",
	})

	fmt.Println("🌱 Initial bank accounts seeded (ACC1001: $5000.00, ACC1002: $1200.00).")
}
