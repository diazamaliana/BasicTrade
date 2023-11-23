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
	"github.com/google/uuid"
)

// CreateVariantRequest represents the request body for creating a new variant.
type CreateVariantRequest struct {
	ProductUUID  string `form:"product_uuid" json:"product_uuid" valid:"required"`
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

func UpdateVariant(c *gin.Context) {
	// Your logic to update a variant
}

func DeleteVariant(c *gin.Context) {
	// Your logic to delete a variant
}

func GetVariantDetail(c *gin.Context) {
	// Your logic to get variant details
}
