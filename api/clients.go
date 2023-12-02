package api

import (
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func clientsRoutes(r *gin.Engine, clientController *controllers.ClientController, cartController *controllers.CartController, db *gorm.DB) {

	clients := r.Group("/clientes")
	clients.POST("", clientController.CreateClient)
	clients.GET("", clientController.GetClients)
	clients.GET("/:client_id", clientController.GetClient)
	clients.PUT("/:client_id", clientController.UpdateClient)
	clients.DELETE("/:client_id", clientController.DeleteClient)

	clients.GET("/:client_id/carrinho", cartController.GetCart)

}
