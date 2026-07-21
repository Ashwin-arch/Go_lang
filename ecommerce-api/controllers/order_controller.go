package controllers

import (
	"fmt"
	"net/http"

	"ecommerce-api/database"
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Checkout converts current cart items into a completed purchase Order atomically
func Checkout(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var cartItems []models.CartItem
	if err := database.DB.Preload("Product").Where("user_id = ?", user.ID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Your shopping cart is empty"})
		return
	}

	var createdOrder models.Order

	// Execute Atomic Database Transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []models.OrderItem

		for _, item := range cartItems {
			var product models.Product
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("product ID %d not found", item.ProductID)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product '%s' (available: %d, requested: %d)", product.Name, product.Stock, item.Quantity)
			}

			// Deduct product stock
			product.Stock -= item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return err
			}

			itemTotal := product.Price * float64(item.Quantity)
			totalAmount += itemTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: product.ID,
				Quantity:  item.Quantity,
				UnitPrice: product.Price,
			})
		}

		// Create Order
		order := models.Order{
			UserID:      user.ID,
			TotalAmount: totalAmount,
			Status:      models.StatusPaid,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		// Attach OrderID to OrderItems and bulk insert
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}

		// Clear user cart
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		createdOrder = order
		createdOrder.OrderItems = orderItems
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Checkout completed successfully! Order placed.",
		"order":   createdOrder,
	})
}

// GetOrders returns all orders placed by the user
func GetOrders(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var orders []models.Order
	if err := database.DB.Preload("OrderItems.Product").Where("user_id = ?", user.ID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
