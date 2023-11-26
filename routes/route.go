package routes

import (
	"basictrade/controllers"
	"basictrade/middleware"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Product routes
	product := router.Group("/products")
	{
		// Middleware
		product.Use(middleware.AuthMiddleware())

		product.GET("", controllers.GetAllProducts)
		product.POST("", controllers.CreateProduct)
		product.PUT("/:productUUID", middleware.ValidateProductAuthorization(),controllers.UpdateProduct)
		product.DELETE("/:productUUID",middleware.ValidateProductAuthorization(), controllers.DeleteProduct)
		product.GET("/:productUUID", controllers.GetProductDetail)

		// Variant routes
		product.GET("/variants", controllers.GetAllVariants)
		product.POST("/variants", controllers.CreateVariant)
		product.PUT("/variants/:variantUUID", middleware.ValidateVariantAuthorization(),controllers.UpdateVariant)
		product.DELETE("/variants/:variantUUID", middleware.ValidateVariantAuthorization(),controllers.DeleteVariant)
		product.GET("/variants/:variantUUID", controllers.GetVariantDetail)
	}

	return router
}
