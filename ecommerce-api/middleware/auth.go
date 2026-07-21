package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"ecommerce-api/database"
	"ecommerce-api/models"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware extracts authentication user ID from Authorization or X-User-ID headers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uint

		// 1. Check Header X-User-ID for testing ease
		userIDHeader := c.GetHeader("X-User-ID")
		if userIDHeader != "" {
			id, err := strconv.ParseUint(userIDHeader, 10, 32)
			if err == nil {
				userID = uint(id)
			}
		}

		// 2. Check Authorization Bearer header
		if userID == 0 {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				id, err := strconv.ParseUint(tokenStr, 10, 32)
				if err == nil {
					userID = uint(id)
				}
			}
		}

		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Please provide valid X-User-ID or Bearer token header."})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User account not found."})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// AdminMiddleware verifies that the authenticated user has Admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated request."})
			c.Abort()
			return
		}

		user := userAny.(models.User)
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin privileges required."})
			c.Abort()
			return
		}

		c.Next()
	}
}
