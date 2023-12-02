package api

import (
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func productsRoutes(r *gin.Engine, productController *controllers.ProductsController, db *gorm.DB) {

	products := r.Group("/produtos")

	products.POST("/", productController.CreateProduct)
	products.GET("/", productController.GetProducts)
	products.GET("/:product_id", productController.GetProduct)
	products.PUT("/:product_id", productController.UpdateProduct)
	products.DELETE("/:product_id", productController.DeleteProduct)

	products.GET("/:product_id/orderswith", productController.GetOrderWithProduct)
}
