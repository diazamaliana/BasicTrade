package middleware

import (
	"basictrade/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the JWT token in the request headers.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify the token
		claims, err := utils.VerifyToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized", 
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Set the claims in the context for further use
		c.Set("adminData", claims)

		// Continue with the next middleware or route handler
		c.Next()
	}
}
