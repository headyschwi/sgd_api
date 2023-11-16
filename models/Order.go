package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Order struct {
	gorm.Model
	ClientID    uint            `json:"client_id" gorm:"not null"`
	OrderItems  []OrderItem     `json:"order_items"`
	TotalPrice  decimal.Decimal `json:"total_price" gorm:"not null"`
	OrderStatus string          `json:"order_status" gorm:"not null"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint            `json:"order_id" gorm:"not null"`
	ProductID uint            `json:"product_id" gorm:"not null"`
	Amount    int64           `json:"quantity" gorm:"not null"`
	Price     decimal.Decimal `json:"price" gorm:"not null"`
}
