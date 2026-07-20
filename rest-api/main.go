package main

import (
	"fmt"
	"log"
	"net/http"

	"rest-api/handlers"
	"rest-api/models"
)

func main() {
	// Initialize in-memory data store
	store := models.NewBookStore()

	// Initialize HTTP handlers with the data store
	bookHandler := handlers.NewBookHandler(store)

	// Create a new HTTP request multiplexer (router)
	mux := http.NewServeMux()

	// Register routes using Go 1.22+ HTTP routing patterns
	mux.HandleFunc("GET /api/books", bookHandler.GetBooks)
	mux.HandleFunc("GET /api/books/{id}", bookHandler.GetBookByID)
	mux.HandleFunc("POST /api/books", bookHandler.CreateBook)
	mux.HandleFunc("PUT /api/books/{id}", bookHandler.UpdateBook)
	mux.HandleFunc("DELETE /api/books/{id}", bookHandler.DeleteBook)

	// Root welcome endpoint
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "Welcome to the Go REST API! Explore endpoints at /api/books")
	})

	port := ":8080"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /api/books       - Fetch all books")
	fmt.Println("  GET    /api/books/{id}  - Fetch book by ID")
	fmt.Println("  POST   /api/books       - Create a new book")
	fmt.Println("  PUT    /api/books/{id}  - Update a book by ID")
	fmt.Println("  DELETE /api/books/{id}  - Delete a book by ID")

	// Start the HTTP server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
