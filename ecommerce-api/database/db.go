package database

import (
	"fmt"
	"log"

	"ecommerce-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes SQLite GORM database connection and performs auto-migrations
func InitDB(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("📦 SQLite Database connected successfully.")

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	SeedData(DB)
	return DB
}

// SeedData inserts initial demo data for testing
func SeedData(db *gorm.DB) {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		return // Data already seeded
	}

	// Seed Users
	customerHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	customer := models.User{Name: "Demo Customer", Email: "customer@example.com", PasswordHash: string(customerHash), Role: "customer"}
	admin := models.User{Name: "Store Admin", Email: "admin@example.com", PasswordHash: string(adminHash), Role: "admin"}
	db.Create(&customer)
	db.Create(&admin)

	// Seed Categories
	electronics := models.Category{Name: "Electronics", Description: "Gadgets, devices and consumer electronics"}
	books := models.Category{Name: "Books", Description: "Technical, fiction and non-fiction books"}
	db.Create(&electronics)
	db.Create(&books)

	// Seed Products
	db.Create(&models.Product{
		Name:        "Wireless Headphones",
		Description: "High fidelity noise cancelling wireless bluetooth headphones.",
		Price:       199.99,
		Stock:       25,
		CategoryID:  electronics.ID,
	})
	db.Create(&models.Product{
		Name:        "Mechanical Gaming Keyboard",
		Description: "RGB back-lit tactile mechanical switches keyboard.",
		Price:       129.50,
		Stock:       15,
		CategoryID:  electronics.ID,
	})
	db.Create(&models.Product{
		Name:        "The Go Programming Language",
		Description: "Comprehensive guide to writing Go idiomatic code by Alan Donovan.",
		Price:       39.99,
		Stock:       50,
		CategoryID:  books.ID,
	})

	fmt.Println("🌱 Initial demo categories, products, and users seeded.")
}
