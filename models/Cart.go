package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Cart struct {
	gorm.Model
	ClientID   uint            `gorm:"not null"`
	CartItems  []CartItem      `gorm:"foreignkey:CartID;association_foreignkey:ID;delete:cascade"`
	TotalPrice decimal.Decimal `gorm:"not null"`
}

type CartItem struct {
	gorm.Model
	CartID    uint            `gorm:"not null" json:"cart_id"`
	ProductID uint            `gorm:"not null" json:"product_id"`
	Amount    int64           `gorm:"not null" json:"amount"`
	Price     decimal.Decimal `gorm:"not null" json:"price"`
}
