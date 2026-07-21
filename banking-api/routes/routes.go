package routes

import (
	"net/http"

	"banking-api/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter initializes Echo router instance with CORS, Logger, HTML UI, and API endpoints
func SetupRouter() *echo.Echo {
	e := echo.New()

	// Built-in Echo Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Static assets and Web UI HTML page
	e.Static("/static", "static")
	e.File("/", "templates/index.html")

	api := e.Group("/api/v1")

	// Accounts & Banking Ledger Routes
	api.POST("/accounts", controllers.CreateAccount)
	api.GET("/accounts", controllers.GetAccounts)
	api.GET("/accounts/:id", controllers.GetAccountByID)
	api.POST("/accounts/:id/deposit", controllers.Deposit)
	api.POST("/accounts/:id/withdraw", controllers.Withdraw)
	api.POST("/transfer", controllers.Transfer)
	api.GET("/accounts/:id/transactions", controllers.GetAccountTransactions)

	return e
}
