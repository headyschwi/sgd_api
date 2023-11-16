// Package main Swagger API documentation.
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta

package main

import (
	"sgd_api/api"
	"sgd_api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Run the API server
//
// swagger:operation POST /api/run operationID
//
// ---
// responses:
//
//	'200':
//	  description: Server started successfully
func main() {
	db := database_conection()
	db.Migrator().DropTable(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	db.Migrator().AutoMigrate(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	api.Run(8080, db)
}

// Establish a database connection
//
// swagger:operation GET /api/database_conection operationID
//
// ---
// responses:
//
//	'200':
//	  description: Database connected successfully
func database_conection() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/sgd_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	print("Database connected successfully!\n")
	return db
}
