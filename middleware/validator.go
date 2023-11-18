package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		validate = validator.New()
		c.Set("validate", validate)
		c.Next()
	}
}
