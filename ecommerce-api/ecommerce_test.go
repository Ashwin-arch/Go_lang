package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce-api/database"
	"ecommerce-api/routes"
)

func setupTestApp() http.Handler {
	database.InitDB(":memory:")
	return routes.SetupRouter()
}

func TestEcommerceAPI(t *testing.T) {
	app := setupTestApp()

	// 1. Test Fetch Products
	t.Run("Get Products", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/products", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})

	// 2. Test Add Item to Cart
	t.Run("Add to Cart", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"productId": 1,
			"quantity":  2,
		})
		req := httptest.NewRequest("POST", "/api/v1/cart/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", "1") // Customer User ID
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201 Created, got %d (body: %s)", w.Code, w.Body.String())
		}
	})

	// 3. Test View Cart
	t.Run("Get Cart", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/cart", nil)
		req.Header.Set("X-User-ID", "1")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})

	// 4. Test Atomic Checkout Transaction
	t.Run("Atomic Checkout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/orders/checkout", nil)
		req.Header.Set("X-User-ID", "1")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201 Created for checkout, got %d (body: %s)", w.Code, w.Body.String())
		}
	})
}
