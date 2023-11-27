// controllers/product_controller.go

package controllers

import (
	"basictrade/models"
	"basictrade/utils"
	"math"
	"mime/multipart"

	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ProductCreateRequest represents the request body for creating a new product.
type ProductCreateRequest struct {
	ProductName     string `form:"product_name" json:"product_name" valid:"required"`
	ImageURL string `form:"image_url" json:"image_url"`
	Image  *multipart.FileHeader `form:"file"`
}

// ProductDetailResponse represents the response structure for product details.
type ProductDetailResponse struct {
    ID         uint   `json:"id"`
    UUID       string `json:"uuid"`
    ProductName string `json:"product_name"`
    ImageURL    string `json:"image_url"`
}

// GetAllProducts retrieves all products from the database with pagination and search.
func GetAllProducts(c *gin.Context) {
	db := utils.GetDB()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	productName := strings.TrimSpace(c.Query("productName"))

	// Pagination logic
	offset := (page - 1) * pageSize

	// Build the query
	query := db.Model(&models.Product{}).Preload("Variants")

	// Apply search filter if name is provided
	if productName!= "" {
		query = query.Where("product_name LIKE ?", "%"+productName+"%")
	}

	// Fetch total count of products
    var totalItems int64
    if err := query.Count(&totalItems).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total items"})
        return
    }

	// Fetch products with pagination
	var products []models.Product
	if err := query.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Calculate total pages
    totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	c.JSON(http.StatusOK, gin.H{"products": products, "totalItems": totalItems, "totalPages": totalPages})
}

// CreateProduct creates a new product.
func CreateProduct(c *gin.Context) {
	// Access claims from the context
	adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db := utils.GetDB()
	contentType := utils.GetContentType(c)

	var createReq ProductCreateRequest
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

	// Generate a unique filename using UUID
	fileName := utils.RemoveExtension(createReq.Image.Filename)

	// Upload the file to Cloudinary
	imageURL, err := utils.UploadFile(createReq.Image, fileName)
	if err != nil {
	   c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to upload file!"})
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

	// Use adminUUID when creating a new product
	newProduct := models.Product{
		ProductName: createReq.ProductName,
		ImageURL:    imageURL,
		AdminUUID:   adminUUID,  // Use the extracted admin UUID
	}

	if err := db.Debug().Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create product", "messages": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": newProduct})
}

// UpdateProduct updates the details of a product.
func UpdateProduct(c *gin.Context) {
	// Access claims from the context
	adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db := utils.GetDB()
	contentType := utils.GetContentType(c)

	// Extract product UUID from the request URL
	productUUIDStr := c.Param("productUUID")

	// Convert product UUID string to uuid.UUID
	productUUID, err := uuid.Parse(productUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product UUID format"})
		return
	}

	var updateReq ProductCreateRequest
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

	// Extract admin UUID from claims
	adminUUIDStr := adminData["adminUUID"].(string)

	// Convert admin UUID string to uuid.UUID
	adminUUID, err := uuid.Parse(adminUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "messages": "Invalid admin UUID format"})
		return
	}

	// Check if the product exists
	var existingProduct models.Product
	if err := db.Where("uuid = ?", productUUID).First(&existingProduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Product not found"})
		return
	}

	// Check if the admin owns the product
	if existingProduct.AdminUUID != adminUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "messages": "You don't have permission to update this product."})
		return
	}

	// Check if the user uploaded a file
    if updateReq.Image != nil {
        // Generate a unique filename using UUID
        fileName := utils.RemoveExtension(updateReq.Image.Filename)

        // Upload the file to Cloudinary
        imageURL, err := utils.UploadFile(updateReq.Image, fileName)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to upload file!"})
            return
        }

        // Update product details
        existingProduct.ImageURL = imageURL
    } else if updateReq.ImageURL != "" {
        // Update product details with the provided image URL
        existingProduct.ImageURL = updateReq.ImageURL
    }

    // Update other product details
    existingProduct.ProductName = updateReq.ProductName

	// Save the updated product details
	if err := db.Save(&existingProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": existingProduct})
}

// DeleteProduct deletes a product.
func DeleteProduct(c *gin.Context) {
	// Access claims from the context
	adminData, exists := c.MustGet("adminData").(jwt5.MapClaims)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	db := utils.GetDB()

	// Extract product UUID from the request URL
	productUUIDStr := c.Param("productUUID")

	// Convert product UUID string to uuid.UUID
	productUUID, err := uuid.Parse(productUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product UUID format"})
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

	// Check if the product exists
	var existingProduct models.Product
	if err := db.Where("uuid = ?", productUUID).Preload("Variants").First(&existingProduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error(), "messages": "Product not found"})
		return
	}

	// Check if the admin owns the product
	if existingProduct.AdminUUID != adminUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error(), "messages": "You don't have permission to delete this product"})
		return
	}

	// Check if the product has associated variants
	if len(existingProduct.Variants) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete product with associated variants", "messages": "Please delete the variants first"})
		return
	}

	// Delete the product
	if err := db.Delete(&existingProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "messages": "Failed to delete product",})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetProductDetail retrieves details of a specific product by UUID.
func GetProductDetail(c *gin.Context) {
    db := utils.GetDB()

    // Extract product UUID from the request URL
    productUUIDStr := c.Param("productUUID")

    // Convert product UUID string to uuid.UUID
    productUUID, err := uuid.Parse(productUUIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product UUID format"})
        return
    }

    // Fetch product details from the database
    var product models.Product
    if err := db.Where("uuid = ?", productUUID).Preload("Variants").First(&product).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error(),"messages": "Product not found"})
        return
    }

    // Create a response struct without sensitive information
    response := ProductDetailResponse{
        ID:         product.ID,
        UUID:       productUUIDStr,
        ProductName: product.ProductName,
        ImageURL:    product.ImageURL,
    }

    c.JSON(http.StatusOK, gin.H{"product": response})
}
