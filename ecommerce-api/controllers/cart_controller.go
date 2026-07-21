package controllers

import (
	"net/http"

	"ecommerce-api/database"
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
)

// GetCart returns items in the authenticated user's cart
func GetCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var items []models.CartItem
	if err := database.DB.Preload("Product.Category").Where("user_id = ?", user.ID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var subtotal float64
	for _, item := range items {
		subtotal += item.Product.Price * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":     items,
		"subtotal": subtotal,
	})
}

// AddToCart adds or updates product quantity in the user's cart
func AddToCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var input struct {
		ProductID uint `json:"productId" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock available"})
		return
	}

	var existingItem models.CartItem
	err := database.DB.Where("user_id = ? AND product_id = ?", user.ID, input.ProductID).First(&existingItem).Error
	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += input.Quantity
		database.DB.Save(&existingItem)
		c.JSON(http.StatusOK, gin.H{"message": "Cart item quantity updated", "item": existingItem})
		return
	}

	newItem := models.CartItem{
		UserID:    user.ID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	database.DB.Create(&newItem)
	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart", "item": newItem})
}

// RemoveFromCart removes an item from cart
func RemoveFromCart(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	itemID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", itemID, user.ID).Delete(&models.CartItem{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
