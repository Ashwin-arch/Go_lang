package models

import (
	"time"
)

// Account represents a bank customer's checking or savings account
type Account struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AccountNumber string    `gorm:"uniqueIndex;not null" json:"accountNumber"`
	AccountHolder string    `gorm:"not null" json:"accountHolder"`
	Balance       float64   `gorm:"not null;default:0.00;check:balance >= 0" json:"balance"`
	Currency      string    `gorm:"default:'USD'" json:"currency"`
	Status        string    `gorm:"default:'ACTIVE'" json:"status"` // 'ACTIVE' or 'FROZEN'
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TransactionType defines the nature of financial transfer
type TransactionType string

const (
	TxTypeDeposit    TransactionType = "DEPOSIT"
	TxTypeWithdrawal TransactionType = "WITHDRAWAL"
	TxTypeTransfer   TransactionType = "TRANSFER"
)

// Transaction represents an immutable ledger entry of funds movement
type Transaction struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	ReferenceNo   string          `gorm:"uniqueIndex;not null" json:"referenceNo"`
	FromAccountID *uint           `gorm:"index" json:"fromAccountId,omitempty"`
	ToAccountID   *uint           `gorm:"index" json:"toAccountId,omitempty"`
	Amount        float64         `gorm:"not null;check:amount > 0" json:"amount"`
	Type          TransactionType `gorm:"not null" json:"type"`
	Status        string          `gorm:"default:'SUCCESS'" json:"status"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"createdAt"`
}

// AuditLog records security & compliance events
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID uint      `gorm:"index" json:"accountId"`
	Action    string    `gorm:"not null" json:"action"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}
