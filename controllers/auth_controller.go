package controllers

import (
	"net/http"

	"basictrade/models"
	"basictrade/utils"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// AdminRequest represents the request body for admin registration.
type AdminRequest struct {
	Name     string `form:"name" json:"name" valid:"required"`
	Email    string `form:"email" json:"email" valid:"email,required"`
	Password string `form:"password" json:"password" valid:"required"`
}

var (
	appJSON = "application/json"
)

// RegisterAdmin handles the registration of a new admin.
func Register(c *gin.Context) {
	db := utils.GetDB()
	contentType := utils.GetContentType(c)
	var adminReq AdminRequest

	// Parse the request body, supports both JSON and form data
	if contentType == appJSON {
		if err := c.ShouldBindJSON(&adminReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBind(&adminReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Validate the input using govalidator
	if _, err := govalidator.ValidateStruct(adminReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Admin instance
	newAdmin := models.Admin{
		Name:     adminReq.Name,
		Email:    adminReq.Email,
		Password: adminReq.Password,
	}

	// Hash the admin's password
	hashedPassword, err := utils.HashPassword(adminReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password!"})
		return
	}
	newAdmin.Password = hashedPassword

	// Save the admin to the database
	if err := db.Debug().Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to register admin!",
			"message": err.Error(),
		})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newAdmin,
		"message": "Admin registered successfully!",
	})

}

// Login handles the login of an admin.
func Login(c *gin.Context) {
	db := utils.GetDB()
	contentType := utils.GetContentType(c)
	var loginReq AdminRequest

	// Parse the request body, supports both JSON and form data
	if contentType == appJSON {
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := c.ShouldBind(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Find the admin with the provided email
	var admin models.Admin
	if err := db.Debug().Where("email = ?", loginReq.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "message": err.Error(),})
		return
	}

	// Verify the password
	if err := utils.VerifyPassword(admin.Password, loginReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "message": err.Error(),})
		return
	}

	// Create a JWT token
	token, err := utils.GenerateToken(admin.UUID, admin.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JWT token", "message": err.Error(),})
		return
	}

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
