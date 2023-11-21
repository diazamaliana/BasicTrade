// controllers/product_controller.go

package controllers

import (
	"basictrade/models"
	"basictrade/utils"
	"fmt"

	// "fmt"

	// "log"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

// ProductCreateRequest represents the request body for creating a new product.
type ProductCreateRequest struct {
	ProductName     string `form:"product_name" json:"product_name" valid:"required"`
	ImageURL string `form:"image_url" json:"image_url" valid:"required"`
}

// GetAllProducts retrieves all products from the database with pagination and search.
func GetAllProducts(c *gin.Context) {
	db := utils.GetDB()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := strings.TrimSpace(c.Query("name"))

	// Pagination logic
	offset := (page - 1) * pageSize

	// Build the query
	query := db.Model(&models.Product{})

	// Apply search filter if name is provided
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Fetch products with pagination
	var products []models.Product
	if err := query.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
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

	// Extract admin UUID from claims
	adminUUIDStr := adminData["adminUUID"].(string)

	// Convert admin UUID string to uuid.UUID
	adminUUID, err := uuid.Parse(adminUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid admin UUID format"})
		return
	}
	fmt.Println("Admin UUID:", adminUUID)  // Add this line for debugging

	// Use adminUUID when creating a new product
	newProduct := models.Product{
		ProductName: createReq.ProductName,
		ImageURL:    createReq.ImageURL,
		AdminUUID:   adminUUID,  // Use the extracted admin UUID
	}

	if err := db.Debug().Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create product", "messages": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": newProduct})
}

func UpdateProduct(c *gin.Context) {
	// Your logic to update a product
}

func DeleteProduct(c *gin.Context) {
	// Your logic to delete a product
}

func GetProductDetail(c *gin.Context) {
	// Your logic to get product details
}
