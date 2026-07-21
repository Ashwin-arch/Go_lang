package controllers

import (
	"net/http"

	"ecommerce-api/database"
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
)

// GetCategories returns all product categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory creates a new category (Admin only)
func CreateCategory(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name, Description: input.Description}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"category": category})
}

// GetProducts returns products with optional filtering by categoryId and query search
func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB.Preload("Category")

	categoryID := c.Query("categoryId")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	search := c.Query("q")
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProductByID fetches a single product
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

// CreateProduct creates a new product (Admin only)
func CreateProduct(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" binding:"gte=0"`
		CategoryID  uint    `json:"categoryId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"product": product})
}

// UpdateProduct updates product details (Admin only)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&product).Updates(input)
	c.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": product})
}

// DeleteProduct removes a product (Admin only)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
