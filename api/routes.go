package api

import (
	"fmt"
	"sgd_api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Run(port int, db *gorm.DB) {
	p := fmt.Sprintf(":%d", port)

	cartController := controllers.NewCartController(db)
	clientController := controllers.NewClientController(db)
	productController := controllers.NewProductController(db)
	orderController := controllers.NewOrderController(db)

	r := gin.Default()

	clientsRoutes(r, clientController, cartController, db)
	cartRoutes(r, cartController, db)
	productsRoutes(r, productController, db)
	ordersRoutes(r, orderController, db)

	r.Run(p)
}
