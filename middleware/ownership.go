package middleware

import (
	"basictrade/models"
	"basictrade/utils"
	"net/http"

	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// ValidateProductAuthorization is a middleware function to check admin authorization.
func ValidateProductAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Access claims from the context
		adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract admin UUID from claims
		adminUUIDStr := adminData["adminUUID"].(string)

		// Convert admin UUID string to uuid.UUID
		adminUUID, err := uuid.Parse(adminUUIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "messages": "Invalid admin UUID format"})
			c.Abort()
			return
		}

		// Extract product UUID from the request URL
		productUUIDStr := c.Param("productUUID")

		// Convert product UUID string to uuid.UUID
		productUUID, err := uuid.Parse(productUUIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product UUID format"})
			c.Abort()
			return
		}

		// Check if the admin owns the product
		db := utils.GetDB()
		var existingProduct models.Product
		if err := db.Where("uuid = ?", productUUID).First(&existingProduct).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Product not found"})
			c.Abort()
			return
		}

		if existingProduct.AdminUUID != adminUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to perform this operation"})
			c.Abort()
			return
		}

		// Set the product in the context for later use
		c.Set("product", existingProduct)

		// Continue with the next middleware or the main handler
		c.Next()
	}
}


func ValidateVariantAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
        // Access claims from the context
        adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Extract admin UUID from claims
        adminUUIDStr := adminData["adminUUID"].(string)

        // Convert admin UUID string to uuid.UUID
        adminUUID, err := uuid.Parse(adminUUIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "messages": "Invalid admin UUID format"})
            c.Abort()
            return
        }

        // Extract variant UUID from the request URL
        variantUUIDStr := c.Param("variantUUID")

        // Convert variant UUID string to uuid.UUID
        variantUUID, err := uuid.Parse(variantUUIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant UUID format"})
            c.Abort()
            return
        }

        // Check if the variant's associated product is owned by the admin
        db := utils.GetDB()
        var existingVariant models.Variant
        if err := db.Where("uuid = ?", variantUUID).First(&existingVariant).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Variant not found"})
            c.Abort()
            return
        }

        var existingProduct models.Product
        if err := db.Model(&existingVariant).Association("Product").Find(&existingProduct); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get associated product"})
            c.Abort()
            return
        }

        if existingProduct.AdminUUID != adminUUID {
            c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to perform this operation"})
            c.Abort()
            return
        }

        // Continue with the next middleware or the main handler
        c.Next()
    }
}
