package api

import (
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func cartRoutes(r *gin.Engine, cartController *controllers.CartController, db *gorm.DB) {

	cart := r.Group("/carrinho")
	cart.GET("/", cartController.GetCart)
	cart.POST("/adicionar/", cartController.AddToCart)
	cart.PATCH("/atualizar/:product_id", cartController.UpdateCartItem)
	cart.DELETE("/remover/:product_id", cartController.RemoveFromCart)

	cart.POST("/finalizar_compra", cartController.Checkout)
}
