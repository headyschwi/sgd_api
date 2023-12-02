package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

type Product struct {
	gorm.Model
	Name         string          `json:"name" gorm:"not null; unique"`
	Description  string          `json:"description"`
	Category     string          `json:"category"`
	Manufacturer string          `json:"manufacturer"`
	Stock        int64           `json:"stock" gorm:"not null"`
	Price        decimal.Decimal `json:"price" gorm:"not null"`
	Weight       float64         `json:"weight"`
	Image_url    string          `json:"image_url"`
}
