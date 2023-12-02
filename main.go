package main

import (
	"fmt"
	"sgd_api/api"
	"sgd_api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := database_connection()
	db.Migrator().DropTable(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	db.Migrator().AutoMigrate(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	api.Run(7777, db)
}

func database_connection() *gorm.DB {
	username := "root"
	password := "root"
	host := "localhost"
	port := "3306"
	dbname := "sgd_api"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	print("Database connected successfully!\n")
	return db
}
