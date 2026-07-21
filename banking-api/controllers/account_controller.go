package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"banking-api/database"
	"banking-api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func generateRefNo() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	return "TXN-" + hex.EncodeToString(bytes)
}

func generateAccNo() string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	return "ACC" + hex.EncodeToString(bytes)
}

// CreateAccount creates a new bank account
func CreateAccount(c echo.Context) error {
	var input struct {
		AccountHolder  string  `json:"accountHolder"`
		InitialDeposit float64 `json:"initialDeposit"`
		Currency       string  `json:"currency"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if input.AccountHolder == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Account holder name is required"})
	}
	if input.InitialDeposit < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Initial deposit cannot be negative"})
	}

	currency := "USD"
	if input.Currency != "" {
		currency = input.Currency
	}

	account := models.Account{
		AccountNumber: generateAccNo(),
		AccountHolder: input.AccountHolder,
		Balance:       input.InitialDeposit,
		Currency:      currency,
		Status:        "ACTIVE",
	}

	if err := database.DB.Create(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create account"})
	}

	// Create Audit Log
	database.DB.Create(&models.AuditLog{
		AccountID: account.ID,
		Action:    "ACCOUNT_CREATED",
		Details:   fmt.Sprintf("Created account %s for %s with initial deposit $%.2f", account.AccountNumber, account.AccountHolder, account.Balance),
		Timestamp: time.Now(),
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Account created successfully",
		"account": account,
	})
}

// GetAccounts lists all registered bank accounts
func GetAccounts(c echo.Context) error {
	var accounts []models.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"accounts": accounts})
}

// GetAccountByID returns details for a specific bank account
func GetAccountByID(c echo.Context) error {
	id := c.Param("id")
	var account models.Account
	if err := database.DB.First(&account, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"account": account})
}

// Deposit credits funds to a bank account
func Deposit(c echo.Context) error {
	id := c.Param("id")
	var input struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := c.Bind(&input); err != nil || input.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Deposit amount must be greater than zero"})
	}

	var updatedAcc models.Account
	var txRecord models.Transaction

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, id).Error; err != nil {
			return fmt.Errorf("account not found")
		}

		if account.Status != "ACTIVE" {
			return fmt.Errorf("account is %s", account.Status)
		}

		account.Balance += input.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		desc := input.Description
		if desc == "" {
			desc = "Cash Deposit"
		}

		txRecord = models.Transaction{
			ReferenceNo: generateRefNo(),
			ToAccountID: &account.ID,
			Amount:      input.Amount,
			Type:        models.TxTypeDeposit,
			Status:      "SUCCESS",
			Description: desc,
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		tx.Create(&models.AuditLog{
			AccountID: account.ID,
			Action:    "DEPOSIT",
			Details:   fmt.Sprintf("Deposited $%.2f into account %s", input.Amount, account.AccountNumber),
			Timestamp: time.Now(),
		})

		updatedAcc = account
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Deposit successful",
		"account":     updatedAcc,
		"transaction": txRecord,
	})
}

// Withdraw debits funds from a bank account with balance validation
func Withdraw(c echo.Context) error {
	id := c.Param("id")
	var input struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := c.Bind(&input); err != nil || input.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Withdrawal amount must be greater than zero"})
	}

	var updatedAcc models.Account
	var txRecord models.Transaction

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, id).Error; err != nil {
			return fmt.Errorf("account not found")
		}

		if account.Status != "ACTIVE" {
			return fmt.Errorf("account is %s", account.Status)
		}

		if account.Balance < input.Amount {
			return fmt.Errorf("insufficient funds. Available balance: $%.2f", account.Balance)
		}

		account.Balance -= input.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		desc := input.Description
		if desc == "" {
			desc = "Cash Withdrawal"
		}

		txRecord = models.Transaction{
			ReferenceNo:   generateRefNo(),
			FromAccountID: &account.ID,
			Amount:        input.Amount,
			Type:          models.TxTypeWithdrawal,
			Status:        "SUCCESS",
			Description:   desc,
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		tx.Create(&models.AuditLog{
			AccountID: account.ID,
			Action:    "WITHDRAWAL",
			Details:   fmt.Sprintf("Withdrew $%.2f from account %s", input.Amount, account.AccountNumber),
			Timestamp: time.Now(),
		})

		updatedAcc = account
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Withdrawal successful",
		"account":     updatedAcc,
		"transaction": txRecord,
	})
}

// Transfer transfers funds between two bank accounts atomically with pessimistic row locking
func Transfer(c echo.Context) error {
	var input struct {
		FromAccountID uint    `json:"fromAccountId"`
		ToAccountID   uint    `json:"toAccountId"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
	}
	if err := c.Bind(&input); err != nil || input.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Transfer amount must be greater than zero"})
	}

	if input.FromAccountID == input.ToAccountID {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Source and destination accounts must be different"})
	}

	var txRecord models.Transaction

	// Atomic ACID Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var fromAcc models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&fromAcc, input.FromAccountID).Error; err != nil {
			return fmt.Errorf("source account ID %d not found", input.FromAccountID)
		}

		if fromAcc.Status != "ACTIVE" {
			return fmt.Errorf("source account is %s", fromAcc.Status)
		}

		if fromAcc.Balance < input.Amount {
			return fmt.Errorf("insufficient balance in account %s (available: $%.2f, requested: $%.2f)", fromAcc.AccountNumber, fromAcc.Balance, input.Amount)
		}

		var toAcc models.Account
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&toAcc, input.ToAccountID).Error; err != nil {
			return fmt.Errorf("destination account ID %d not found", input.ToAccountID)
		}

		if toAcc.Status != "ACTIVE" {
			return fmt.Errorf("destination account is %s", toAcc.Status)
		}

		// Perform transfer adjustments
		fromAcc.Balance -= input.Amount
		toAcc.Balance += input.Amount

		if err := tx.Save(&fromAcc).Error; err != nil {
			return err
		}
		if err := tx.Save(&toAcc).Error; err != nil {
			return err
		}

		desc := input.Description
		if desc == "" {
			desc = fmt.Sprintf("Fund transfer from %s to %s", fromAcc.AccountNumber, toAcc.AccountNumber)
		}

		txRecord = models.Transaction{
			ReferenceNo:   generateRefNo(),
			FromAccountID: &fromAcc.ID,
			ToAccountID:   &toAcc.ID,
			Amount:        input.Amount,
			Type:          models.TxTypeTransfer,
			Status:        "SUCCESS",
			Description:   desc,
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}

		// Create Audit Logs
		tx.Create(&models.AuditLog{
			AccountID: fromAcc.ID,
			Action:    "TRANSFER_OUT",
			Details:   fmt.Sprintf("Transferred $%.2f to %s", input.Amount, toAcc.AccountNumber),
			Timestamp: time.Now(),
		})
		tx.Create(&models.AuditLog{
			AccountID: toAcc.ID,
			Action:    "TRANSFER_IN",
			Details:   fmt.Sprintf("Received $%.2f from %s", input.Amount, fromAcc.AccountNumber),
			Timestamp: time.Now(),
		})

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Transfer completed successfully",
		"transaction": txRecord,
	})
}

// GetAccountTransactions fetches ledger statement for an account
func GetAccountTransactions(c echo.Context) error {
	id := c.Param("id")

	var transactions []models.Transaction
	err := database.DB.Where("from_account_id = ? OR to_account_id = ?", id, id).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"transactions": transactions})
}
