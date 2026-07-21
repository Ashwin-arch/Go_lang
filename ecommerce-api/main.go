package main

import (
	"fmt"
	"log"
	"os"

	"ecommerce-api/database"
	"ecommerce-api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize SQLite Database
	database.InitDB("ecommerce.db")

	// Setup Gin Router
	r := routes.SetupRouter()

	fmt.Println("🚀 Gin E-commerce REST API starting...")
	fmt.Printf("📡 Base URL           : http://localhost:%s/\n", port)
	fmt.Printf("🛒 Products Endpoint  : GET http://localhost:%s/api/v1/products\n", port)
	fmt.Printf("🛍️ Cart Endpoint      : GET http://localhost:%s/api/v1/cart (Requires X-User-ID: 1)\n", port)
	fmt.Printf("💳 Checkout Endpoint  : POST http://localhost:%s/api/v1/orders/checkout (Requires X-User-ID: 1)\n", port)
	fmt.Println("-------------------------------------------------------------------------")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
