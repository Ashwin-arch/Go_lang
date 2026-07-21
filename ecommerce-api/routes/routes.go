package routes

import (
	"net/http"

	"ecommerce-api/controllers"
	"ecommerce-api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes Gin engine with middlewares, HTML Web UI, and API endpoint routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Static assets and Web UI HTML page
	r.Static("/static", "static")
	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	api := r.Group("/api/v1")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Public Product & Category Routes
		api.GET("/categories", controllers.GetCategories)
		api.GET("/products", controllers.GetProducts)
		api.GET("/products/:id", controllers.GetProductByID)

		// Admin Protected Product Routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.POST("/categories", controllers.CreateCategory)
			admin.POST("/products", controllers.CreateProduct)
			admin.PUT("/products/:id", controllers.UpdateProduct)
			admin.DELETE("/products/:id", controllers.DeleteProduct)
		}

		// Authenticated Customer Cart & Order Routes
		customer := api.Group("")
		customer.Use(middleware.AuthMiddleware())
		{
			customer.GET("/cart", controllers.GetCart)
			customer.POST("/cart/items", controllers.AddToCart)
			customer.DELETE("/cart/items/:id", controllers.RemoveFromCart)

			customer.POST("/orders/checkout", controllers.Checkout)
			customer.GET("/orders", controllers.GetOrders)
		}
	}

	return r
}
