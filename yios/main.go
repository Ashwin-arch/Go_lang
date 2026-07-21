package main

import (
	"fmt"
	"os"

	"yios/database"
	"yios/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Initialize SQLite Database & Seed Data
	database.InitDB("yios.db")

	// Initialize Gin Router & Web UI Engine
	r := routes.SetupRouter()

	fmt.Println("🚀 Yios AI Kubernetes Control Plane starting...")
	fmt.Printf("📡 Base Dashboard URL : http://localhost:%s/\n", port)
	fmt.Printf("🤖 Deployments API    : GET http://localhost:%s/api/v1/deployments\n", port)
	fmt.Printf("📊 K8s Pod Metrics    : GET http://localhost:%s/api/v1/deployments/1/metrics\n", port)
	fmt.Println("-------------------------------------------------------------------------")

	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Error starting Yios platform server: %v\n", err)
	}
}
