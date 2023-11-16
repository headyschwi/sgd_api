package models

import "github.com/jinzhu/gorm"

type Client struct {
	gorm.Model
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null;unique"`

	Cart   Cart    `gorm:"foreignkey:ClientID; delete:cascade"`
	Orders []Order `gorm:"foreignkey:ClientID"`
}
