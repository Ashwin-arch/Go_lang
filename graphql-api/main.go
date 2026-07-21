package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"graphql-api/handlers"
	"graphql-api/models"
	"graphql-api/schema"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize thread-safe data store with sample seed data
	store := models.NewStore()

	// Build GraphQL schema with queries & mutations
	gqlSchema, err := schema.BuildSchema(store)
	if err != nil {
		log.Fatalf("Failed to build GraphQL schema: %v", err)
	}

	// Create GraphQL HTTP handler
	gqlHandler := handlers.NewGraphQLHandler(gqlSchema)

	mux := http.NewServeMux()
	mux.Handle("/graphql", gqlHandler)
	mux.Handle("/", gqlHandler)

	addr := ":" + port
	fmt.Println("🚀 Go GraphQL API Server starting...")
	fmt.Printf("📡 GraphQL Endpoint : http://localhost:%s/graphql\n", port)
	fmt.Printf("🎨 GraphiQL IDE UI  : http://localhost:%s/\n", port)
	fmt.Println("-------------------------------------------------------")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server exited unexpectedly: %v", err)
	}
}
