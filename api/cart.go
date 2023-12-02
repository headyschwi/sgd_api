package api

import (
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func cartRoutes(r *gin.Engine, cartController *controllers.CartController, db *gorm.DB) {

	cart := r.Group("/carrinho")

	cart.POST("/adicionar/", cartController.AddToCart)
	cart.PUT("/atualizar/", cartController.UpdateCartItem)
	cart.DELETE("/remover/", cartController.RemoveFromCart)

	cart.POST("/finalizar_compra", cartController.Checkout)
}
