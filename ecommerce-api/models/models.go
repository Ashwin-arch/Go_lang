package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a customer or administrator in the e-commerce system
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"default:'customer'" json:"role"` // 'customer' or 'admin'
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Category represents a product category
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Product represents an item listed for sale
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CategoryID  uint      `gorm:"not null" json:"categoryId"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	ProductID uint      `gorm:"not null" json:"productId"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

// OrderStatus defines the lifecycle status of an order
type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusCancelled OrderStatus = "CANCELLED"
)

// Order represents a completed customer purchase order
type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null;index" json:"userId"`
	User        User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalAmount float64     `gorm:"not null" json:"totalAmount"`
	Status      OrderStatus `gorm:"default:'PAID'" json:"status"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems"`
	CreatedAt   time.Time   `json:"createdAt"`
}

// OrderItem represents an item snapshot in a placed order
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null;index" json:"orderId"`
	ProductID uint    `gorm:"not null" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	UnitPrice float64 `gorm:"not null" json:"unitPrice"`
}
