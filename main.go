package main

import (
	"fmt"
	"os"
	"sgd_api/api"
	"sgd_api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := database_connection()
	db.Migrator().DropTable(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	db.Migrator().AutoMigrate(&models.Client{}, &models.Cart{}, &models.Order{}, &models.Product{}, &models.CartItem{}, &models.OrderItem{})
	api.Run(8080, db)
}

func database_connection() *gorm.DB {
	// Obtém as variáveis de ambiente
	username := os.Getenv("DB_USERNAME") // substitua 'DB_USERNAME' pelo nome da sua variável de ambiente
	password := os.Getenv("DB_PASSWORD") // substitua 'DB_PASSWORD' pelo nome da sua variável de ambiente
	host := os.Getenv("DB_HOST")         // substitua 'DB_HOST' pelo nome da sua variável de ambiente
	port := os.Getenv("DB_PORT")         // substitua 'DB_PORT' pelo nome da sua variável de ambiente
	dbname := os.Getenv("DB_NAME")       // substitua 'DB_NAME' pelo nome da sua variável de ambiente

	// Monta a string de conexão usando as variáveis de ambiente
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)

	// Conecta ao banco de dados
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	print("Database connected successfully!\n")
	return db
}
