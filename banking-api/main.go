package main

import (
	"fmt"
	"os"

	"banking-api/database"
	"banking-api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Initialize SQLite Database
	database.InitDB("banking.db")

	// Initialize Echo Router
	e := routes.SetupRouter()

	fmt.Println("🚀 Echo Banking REST API starting...")
	fmt.Printf("📡 Base URL            : http://localhost:%s/\n", port)
	fmt.Printf("🏦 Accounts Endpoint   : GET http://localhost:%s/api/v1/accounts\n", port)
	fmt.Printf("💸 Transfer Endpoint   : POST http://localhost:%s/api/v1/transfer\n", port)
	fmt.Println("-------------------------------------------------------------------------")

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}
}
