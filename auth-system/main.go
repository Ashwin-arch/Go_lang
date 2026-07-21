package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"auth-system/handlers"
	"auth-system/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize thread-safe AuthStore with bcrypt password hashing
	store := models.NewAuthStore()

	// Initialize Auth HTTP handlers
	authHandler := handlers.NewAuthHandler(store, "templates")

	mux := http.NewServeMux()

	// Static assets handler
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// API Endpoints
	mux.HandleFunc("POST /api/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/me", authHandler.Me)

	// HTML UI Page
	mux.HandleFunc("GET /", authHandler.ServeIndex)

	addr := ":" + port
	fmt.Println("🚀 Go Authentication Server starting...")
	fmt.Printf("🌐 Web Application UI : http://localhost:%s/\n", port)
	fmt.Printf("🔐 Login API          : POST http://localhost:%s/api/login\n", port)
	fmt.Printf("📝 Signup API         : POST http://localhost:%s/api/signup\n", port)
	fmt.Printf("💡 Demo Credentials   : Username: 'demo', Password: 'password123'\n")
	fmt.Println("----------------------------------------------------------------")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server exited unexpectedly: %v", err)
	}
}
