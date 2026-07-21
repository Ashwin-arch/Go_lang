package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"yios/controllers"
)

// SetupRouter initializes Gin engine with CORS, Web UI templates, and API endpoints
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Static assets and Yios React/TS Glassmorphism SPA Dashboard
	r.Static("/static", "static")
	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	api := r.Group("/api/v1")
	{
		// Authentication Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// AI Model Deployment & Kubernetes Routes
		deploy := api.Group("/deployments")
		{
			deploy.GET("", controllers.GetDeployments)
			deploy.POST("", controllers.CreateDeployment)
			deploy.GET("/:id", controllers.GetDeploymentByID)
			deploy.POST("/:id/scale", controllers.ScaleDeployment)
			deploy.POST("/:id/update", controllers.UpdateDeployment)
			deploy.POST("/:id/rollback", controllers.RollbackDeployment)
			deploy.GET("/:id/revisions", controllers.GetRevisions)
			deploy.DELETE("/:id", controllers.DeleteDeployment)

			// K8s Real-Time Operations & Monitoring
			deploy.GET("/:id/metrics", controllers.GetMetrics)
			deploy.GET("/:id/logs", controllers.GetPodLogs)
			deploy.GET("/:id/manifest", controllers.GetManifest)
			deploy.POST("/:id/proxy", controllers.ProxyInference)
			deploy.POST("/:id/chat", controllers.ChatInference)
		}
	}

	return r
}
