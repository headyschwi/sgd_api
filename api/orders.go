package api

import (
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ordersRoutes(r *gin.Engine, ordersController *controllers.OrdersController, db *gorm.DB) {

	orders := r.Group("/pedidos")

	orders.GET("", ordersController.GetOrders)
	orders.PATCH("/:order_id", ordersController.UpdateOrderStatus)
	orders.POST("/pedir_denovo/:order_id", ordersController.OrderAgain)

}
