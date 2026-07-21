package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-system/handlers"
	"auth-system/models"
)

func setupTestServer() (*models.AuthStore, http.Handler) {
	store := models.NewAuthStore()
	authHandler := handlers.NewAuthHandler(store, "templates")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/signup", authHandler.Signup)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/me", authHandler.Me)

	return store, mux
}

func TestAuthSystem(t *testing.T) {
	_, mux := setupTestServer()

	// 1. Test Demo User Login
	t.Run("Login with Demo Account", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"usernameOrEmail": "demo",
			"password":        "password123",
		})
		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		// Check if session_token cookie was set
		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "session_token" {
				sessionCookie = c
				break
			}
		}

		if sessionCookie == nil || sessionCookie.Value == "" {
			t.Fatalf("Expected session_token cookie to be set")
		}

		// 2. Test Get Session User (/api/me)
		meReq := httptest.NewRequest("GET", "/api/me", nil)
		meReq.AddCookie(sessionCookie)
		meW := httptest.NewRecorder()

		mux.ServeHTTP(meW, meReq)

		if meW.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for /api/me, got %d", meW.Code)
		}
	})

	// 3. Test Invalid Credentials
	t.Run("Login with Invalid Password", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"usernameOrEmail": "demo",
			"password":        "wrongpassword",
		})
		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected status 401 Unauthorized, got %d", w.Code)
		}
	})

	// 4. Test New User Registration (Signup)
	t.Run("Register New User", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"username": "alice",
			"email":    "alice@example.com",
			"password": "securepassword123",
			"fullName": "Alice Smith",
		})
		req := httptest.NewRequest("POST", "/api/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201 Created, got %d", w.Code)
		}
	})
}
