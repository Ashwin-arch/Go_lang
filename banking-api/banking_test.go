package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"banking-api/database"
	"banking-api/routes"
)

func setupTestApp() http.Handler {
	database.InitDB(":memory:")
	return routes.SetupRouter()
}

func TestBankingAPI(t *testing.T) {
	app := setupTestApp()

	// 1. Test Fetching Accounts (Pre-seeded: Account 1 with $5000, Account 2 with $1200)
	t.Run("Get Accounts", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/accounts", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})

	// 2. Test Deposit into Account 1 ($500)
	t.Run("Deposit Funds", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"amount":      500.00,
			"description": "Salary Bonus Deposit",
		})
		req := httptest.NewRequest("POST", "/api/v1/accounts/1/deposit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for deposit, got %d (body: %s)", w.Code, w.Body.String())
		}
	})

	// 3. Test Withdrawal from Account 1 ($200)
	t.Run("Withdraw Funds", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"amount":      200.00,
			"description": "ATM Cash Withdrawal",
		})
		req := httptest.NewRequest("POST", "/api/v1/accounts/1/withdraw", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for withdrawal, got %d (body: %s)", w.Code, w.Body.String())
		}
	})

	// 4. Test Inter-Bank Fund Transfer (Account 1 -> Account 2, Amount: $1000)
	t.Run("Atomic Fund Transfer", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"fromAccountId": 1,
			"toAccountId":   2,
			"amount":        1000.00,
			"description":   "Rent Payment Transfer",
		})
		req := httptest.NewRequest("POST", "/api/v1/transfer", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200 for transfer, got %d (body: %s)", w.Code, w.Body.String())
		}
	})

	// 5. Test Insufficient Funds Transfer (Attempting to transfer $100,000)
	t.Run("Insufficient Funds Transfer Fail", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"fromAccountId": 2,
			"toAccountId":   1,
			"amount":        100000.00,
			"description":   "Excessive Transfer",
		})
		req := httptest.NewRequest("POST", "/api/v1/transfer", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400 Bad Request for insufficient funds, got %d", w.Code)
		}
	})
}
