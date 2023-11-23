// controllers/variant_controller.go

package controllers

import (
	"basictrade/models"
	"basictrade/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

// CreateVariantRequest represents the request body for creating a new variant.
type CreateVariantRequest struct {
	ProductUUID  string `form:"product_uuid" json:"product_uuid"`
    VariantName string `form:"variant_name" json:"variant_name" valid:"required"`
    Quantity    uint   `form:"quantity" json:"quantity" valid:"required"`
}

func GetAllVariants(c *gin.Context) {
	db := utils.GetDB()

    // Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	variantName := strings.TrimSpace(c.Query("variantName"))

	// Pagination logic
	offset := (page - 1) * pageSize

	// Build the query
	query := db.Model(&models.Variant{})

	// Apply search filter if name is provided
	if variantName!= "" {
		query = query.Where("variant_name LIKE ?", "%"+variantName+"%")
	}

	// Fetch products with pagination
    var variants []models.Variant
    if err := query.Offset(offset).Limit(pageSize).Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

    c.JSON(http.StatusOK, gin.H{"variants": variants})
}

// CreateVariant creates a new variant for a specific product.
func CreateVariant(c *gin.Context) {
    db := utils.GetDB()
	var createReq CreateVariantRequest
	contentType := utils.GetContentType(c)

    if contentType == appJSON {
		if err := c.ShouldBindJSON(&createReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBind(&createReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if _, err := govalidator.ValidateStruct(createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert product UUID string to uuid.UUID
    productUUID, err := uuid.Parse(createReq.ProductUUID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "messages": "Invalid product UUID format"})
        return
    }

	// Check if the product with the given UUID exists
	var existingProduct models.Product
	if err := db.Where("uuid = ?", productUUID).First(&existingProduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Product not found"})
		return
	}

    // Create a new variant
    newVariant := models.Variant{
        VariantName: createReq.VariantName,
        Quantity:    createReq.Quantity,
        ProductUUID:   productUUID,
    }

    // Save the new variant to the database
    if err := db.Create(&newVariant).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to create variant"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"variant": newVariant})
}

// UpdateVariant updates the details of a variant.
func UpdateVariant(c *gin.Context) {    
    db := utils.GetDB()

    // Extract variant UUID from the request URL
    variantUUIDStr := c.Param("variantUUID")

    // Convert variant UUID string to uuid.UUID
    variantUUID, err := uuid.Parse(variantUUIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid variant UUID format"})
        return
    }

    var updateReq CreateVariantRequest
    contentType := utils.GetContentType(c)

    if contentType == appJSON {
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBind(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if _, err := govalidator.ValidateStruct(updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    // Check if the variant exists
    var existingVariant models.Variant
    if err := db.Where("uuid = ?", variantUUID).First(&existingVariant).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Variant not found"})
        return
    }

    // Check if the admin owns the product associated with the variant
    var existingProduct models.Product
    if err := db.Where("uuid = ?", existingVariant.ProductUUID).First(&existingProduct).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to fetch product"})
        return
    }

    // Access claims from the context
    adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Extract admin UUID from claims
    adminUUIDStr := adminData["adminUUID"].(string)

    // Convert admin UUID string to uuid.UUID
	adminUUID, err := uuid.Parse(adminUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin UUID format"})
		return
	}

    // Check if the admin owns the product
    if existingProduct.AdminUUID != adminUUID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this variant."})
        return
    }

    // Update variant details
    existingVariant.VariantName = updateReq.VariantName
    existingVariant.Quantity = updateReq.Quantity

    // Save the updated variant details
    if err := db.Save(&existingVariant).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to update variant"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"variant": existingVariant})
}

func DeleteVariant(c *gin.Context) {
	// Your logic to delete a variant
}

func GetVariantDetail(c *gin.Context) {
	// Your logic to get variant details
}
