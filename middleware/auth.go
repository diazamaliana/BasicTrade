package middleware

import (
	auth "basictrade/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Your authentication logic here
		// Verify JWT token, check user roles, etc.
		// For simplicity, let's assume authentication is handled in another function
		if !auth.IsAuthenticated(c) {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
